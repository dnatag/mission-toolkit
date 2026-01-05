package plan

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/dnatag/mission-toolkit/internal/logger"
	"github.com/spf13/afero"
)

var (
	validDomains = map[string]bool{
		"security": true, "performance": true, "complex-algo": true, "high-risk": true,
		"cross-cutting": true, "authentication": true, "authorization": true,
		"cryptography": true, "real-time": true, "compliance": true,
	}
	suspiciousKeywords = []string{"rm -rf", "sudo", "chmod 777", "eval", "exec"}
	dangerousCommands  = []string{"rm", "delete", "drop", "truncate", "format"}
)

// ValidationResult represents the result of plan validation
type ValidationResult struct {
	Valid          bool     `json:"valid"`
	Errors         []string `json:"errors,omitempty"`
	Warnings       []string `json:"warnings,omitempty"`
	SecurityIssues []string `json:"security_issues,omitempty"`
	FileValidation []string `json:"file_validation,omitempty"`
	FormatIssues   []string `json:"format_issues,omitempty"`
}

// ValidatorService handles plan validation with comprehensive safety checks
type ValidatorService struct {
	ServiceBase
	log     *logger.Logger
	rootDir string
}

// NewValidatorService creates a new validator service
func NewValidatorService(fs afero.Fs, missionID string, rootDir string) *ValidatorService {
	// Create logger config with same filesystem as the service
	config := logger.DefaultConfig()
	config.Fs = fs

	return &ValidatorService{
		ServiceBase: NewServiceBase(fs, missionID),
		log:         logger.NewWithConfig(missionID, config),
		rootDir:     rootDir,
	}
}

// ValidatePlan performs comprehensive validation of a plan.json file
func (v *ValidatorService) ValidatePlan(planPath string) (*ValidationResult, error) {
	result := &ValidationResult{Valid: true}
	v.log.LogStep("INFO", "validation-start", fmt.Sprintf("Validating plan: %s", planPath))

	if exists, _ := afero.Exists(v.fs, planPath); !exists {
		return v.addError(result, "Plan file does not exist: %s", planPath), nil
	}

	planData, err := afero.ReadFile(v.fs, planPath)
	if err != nil {
		return v.addError(result, "Failed to read plan file: %v", err), nil
	}

	var planSpec PlanSpec
	if err := json.Unmarshal(planData, &planSpec); err != nil {
		result.Valid = false
		result.FormatIssues = append(result.FormatIssues, fmt.Sprintf("Invalid JSON format: %v", err))
		return result, nil
	}

	v.validatePlanFormat(&planSpec, result)
	v.validateFilePaths(&planSpec, result)
	v.validateMissionContent(&planSpec, result)

	// Valid only if no errors AND no security issues
	result.Valid = len(result.Errors) == 0 && len(result.SecurityIssues) == 0

	v.log.LogStep("INFO", "validation-complete",
		fmt.Sprintf("Valid: %t, errors: %d, warnings: %d, security: %d",
			result.Valid, len(result.Errors), len(result.Warnings), len(result.SecurityIssues)))

	return result, nil
}

// Helper function to add errors
func (v *ValidatorService) addError(result *ValidationResult, format string, args ...interface{}) *ValidationResult {
	result.Valid = false
	result.Errors = append(result.Errors, fmt.Sprintf(format, args...))
	return result
}

// validatePlanFormat checks the basic structure and required fields
func (v *ValidatorService) validatePlanFormat(spec *PlanSpec, result *ValidationResult) {
	if spec.Intent == "" {
		result.Errors = append(result.Errors, "Missing required field: intent")
	}
	// Check if we have files in either scope or files field
	allFiles := spec.GetScopeFiles()
	if len(allFiles) == 0 {
		result.Errors = append(result.Errors, "Missing required field: scope (must contain at least one file)")
	}
	if len(spec.Plan) == 0 {
		result.Errors = append(result.Errors, "Missing required field: plan (must contain at least one step)")
	}
	if spec.Verification == "" {
		result.Warnings = append(result.Warnings, "Missing verification command - recommended for quality assurance")
	}

	for _, domain := range spec.Domain {
		if !validDomains[strings.ToLower(domain)] {
			result.Warnings = append(result.Warnings, fmt.Sprintf("Unknown domain: %s", domain))
		}
	}
}

// validateFilePaths performs security and accessibility checks on file paths
func (v *ValidatorService) validateFilePaths(spec *PlanSpec, result *ValidationResult) {
	// Use GetScopeFiles to ensure we validate all files (legacy Scope + new Files)
	for _, filePath := range spec.GetScopeFiles() {
		if strings.Contains(filePath, "..") {
			result.SecurityIssues = append(result.SecurityIssues, fmt.Sprintf("Path traversal detected in: %s", filePath))
			continue
		}

		if !v.isWithinRoot(filePath) {
			result.SecurityIssues = append(result.SecurityIssues, fmt.Sprintf("File outside project root: %s", filePath))
			continue
		}

		v.validateFileAccess(filePath, result)
	}
}

// isWithinRoot checks if file path is within project root
func (v *ValidatorService) isWithinRoot(filePath string) bool {
	absPath := filepath.Join(v.rootDir, filePath)
	cleanPath := filepath.Clean(absPath)
	cleanRoot := filepath.Clean(v.rootDir)

	if !filepath.IsAbs(cleanRoot) {
		if wd, err := os.Getwd(); err == nil {
			cleanRoot = filepath.Join(wd, cleanRoot)
		}
	}
	if !filepath.IsAbs(cleanPath) {
		if wd, err := os.Getwd(); err == nil {
			cleanPath = filepath.Join(wd, cleanPath)
		}
	}

	return strings.HasPrefix(cleanPath, cleanRoot)
}

// validateFileAccess checks if files exist and are accessible
func (v *ValidatorService) validateFileAccess(filePath string, result *ValidationResult) {
	fullPath := filepath.Join(v.rootDir, filePath)

	if exists, _ := afero.Exists(v.fs, fullPath); exists {
		if file, err := v.fs.OpenFile(fullPath, os.O_RDWR, 0644); err != nil {
			result.FileValidation = append(result.FileValidation, fmt.Sprintf("File not writable: %s (%v)", filePath, err))
			// Not writable is an error for modification
			result.Errors = append(result.Errors, fmt.Sprintf("File not writable: %s", filePath))
		} else {
			file.Close()
			result.FileValidation = append(result.FileValidation, fmt.Sprintf("File accessible: %s", filePath))
		}
	} else {
		// File does not exist - check if we can create it
		dir := filepath.Dir(fullPath)
		if err := v.fs.MkdirAll(dir, 0755); err != nil {
			result.FileValidation = append(result.FileValidation, fmt.Sprintf("Cannot create directory for: %s (%v)", filePath, err))
			result.Errors = append(result.Errors, fmt.Sprintf("Cannot create directory for: %s", filePath))
		} else {
			result.FileValidation = append(result.FileValidation, fmt.Sprintf("File can be created: %s", filePath))
			// Add explicit warning for new file creation as per HLD
			result.Warnings = append(result.Warnings, fmt.Sprintf("File does not exist and will be created: %s", filePath))
		}
	}
}

// validateMissionContent performs content-level validation
func (v *ValidatorService) validateMissionContent(spec *PlanSpec, result *ValidationResult) {
	v.checkSuspiciousContent(strings.ToLower(spec.Intent), "intent", result)

	for i, step := range spec.Plan {
		v.checkSuspiciousContent(strings.ToLower(step), fmt.Sprintf("plan step %d", i+1), result)
	}

	if spec.Verification != "" {
		v.checkDangerousCommands(strings.ToLower(spec.Verification), result)
	}
}

// checkSuspiciousContent checks for suspicious keywords
func (v *ValidatorService) checkSuspiciousContent(content, location string, result *ValidationResult) {
	for _, keyword := range suspiciousKeywords {
		if strings.Contains(content, keyword) {
			result.SecurityIssues = append(result.SecurityIssues, fmt.Sprintf("Suspicious keyword in %s: %s", location, keyword))
		}
	}
}

// checkDangerousCommands checks for dangerous verification commands
func (v *ValidatorService) checkDangerousCommands(content string, result *ValidationResult) {
	for _, cmd := range dangerousCommands {
		if strings.Contains(content, cmd) {
			result.SecurityIssues = append(result.SecurityIssues, fmt.Sprintf("Potentially dangerous verification command: %s", cmd))
		}
	}
}

// ToJSON converts validation result to JSON string
func (r *ValidationResult) ToJSON() (string, error) {
	return ToJSON(r)
}

// FormatValidationOutput creates standardized output with next_step guidance
func (v *ValidatorService) FormatValidationOutput(result *ValidationResult) OutputResponse {
	nextStep := "PROCEED to Step 6 (Generate Mission)."
	if !result.Valid {
		nextStep = "FIX errors in plan.json and retry validation."
	}

	data := map[string]interface{}{
		"valid":           result.Valid,
		"errors":          result.Errors,
		"warnings":        result.Warnings,
		"security_issues": result.SecurityIssues,
		"file_validation": result.FileValidation,
	}

	return NewOutputResponse(data, nextStep)
}
