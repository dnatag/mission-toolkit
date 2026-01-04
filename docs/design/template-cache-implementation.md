# Template Cache Implementation

**Status**: Ready for Implementation  
**Date**: 2026-01-03

---

## Problem

Current m.plan workflow loads 8+ templates individually:
- Step 1: 4 templates (intent, clarification, missions/clarification, displays/plan-clarification)
- Step 2: 2 templates (duplication, domain)
- Step 4: 1 template (plan-atomic OR plan-epic)
- Step 6: 1 template (plan-success)

**Issue**: Multiple file reads per execution, no reuse across steps.

---

## Solution

Single cache file with `context.step` keys for efficient single-step loading.

### Cache Structure

```json
{
  "version": "1.0.0",
  "created_at": "2026-01-03T10:30:00Z",
  "templates": {
    "analysis/intent": "<content>",
    "analysis/clarification": "<content>",
    "displays/plan-success": "<content>"
  },
  "mappings": {
    "plan.step1": ["analysis/intent", "analysis/clarification"],
    "plan.step6": ["displays/plan-success"],
    "apply.step5": ["displays/apply-success"]
  }
}
```

### CLI Commands

```bash
# Build cache (run once, or auto-rebuild if missing)
m templates cache build

# Get templates for specific step
m templates cache get --context plan --step step1

# Clear cache
m templates cache clear
```

### Output Format

```json
{
  "context": "plan",
  "step": "step1",
  "templates": {
    "analysis/intent": "<content>",
    "analysis/clarification": "<content>"
  }
}
```

---

## Implementation

### 1. Core Cache Logic

**File**: `internal/templates/cache.go`

```go
package templates

import (
    "encoding/json"
    "fmt"
    "path/filepath"
    "time"
    "github.com/spf13/afero"
)

type TemplateCache struct {
    Version   string              `json:"version"`
    CreatedAt time.Time           `json:"created_at"`
    Templates map[string]string   `json:"templates"`
    Mappings  map[string][]string `json:"mappings"`
}

type StepTemplates struct {
    Context   string            `json:"context"`
    Step      string            `json:"step"`
    Templates map[string]string `json:"templates"`
}

func BuildCache(fs afero.Fs) error {
    cache := &TemplateCache{
        Version:   "1.0.0",
        CreatedAt: time.Now(),
        Templates: make(map[string]string),
        Mappings:  getTemplateMappings(),
    }
    
    baseDir := ".mission/libraries"
    paths := []string{
        "analysis/intent", "analysis/clarification", "analysis/duplication",
        "analysis/domain", "missions/clarification", "displays/plan-atomic",
        "displays/plan-epic", "displays/plan-success", "displays/plan-clarification",
    }
    
    for _, path := range paths {
        content, err := afero.ReadFile(fs, filepath.Join(baseDir, path+".md"))
        if err != nil {
            continue // Skip missing templates
        }
        cache.Templates[path] = string(content)
    }
    
    data, _ := json.MarshalIndent(cache, "", "  ")
    return afero.WriteFile(fs, filepath.Join(baseDir, "cache.json"), data, 0644)
}

func GetStepTemplates(fs afero.Fs, context, step string) (*StepTemplates, error) {
    data, err := afero.ReadFile(fs, ".mission/libraries/cache.json")
    if err != nil {
        BuildCache(fs) // Auto-rebuild
        data, err = afero.ReadFile(fs, ".mission/libraries/cache.json")
        if err != nil {
            return nil, err
        }
    }
    
    var cache TemplateCache
    json.Unmarshal(data, &cache)
    
    key := fmt.Sprintf("%s.%s", context, step)
    paths := cache.Mappings[key]
    
    result := &StepTemplates{
        Context:   context,
        Step:      step,
        Templates: make(map[string]string),
    }
    
    for _, path := range paths {
        result.Templates[path] = cache.Templates[path]
    }
    
    return result, nil
}

func getTemplateMappings() map[string][]string {
    return map[string][]string{
        "plan.step1": {"analysis/intent", "analysis/clarification", "missions/clarification", "displays/plan-clarification"},
        "plan.step2": {"analysis/duplication", "analysis/domain"},
        "plan.step4_track1": {"displays/plan-atomic"},
        "plan.step4_track4": {"displays/plan-epic"},
        "plan.step6": {"displays/plan-success"},
    }
}
```

### 2. CLI Commands

**File**: `cmd/templates.go`

```go
package cmd

import (
    "encoding/json"
    "os"
    "github.com/spf13/cobra"
    "github.com/spf13/afero"
    "your-project/internal/templates"
)

var templatesCmd = &cobra.Command{Use: "templates"}
var cacheCmd = &cobra.Command{Use: "cache"}

var buildCmd = &cobra.Command{
    Use: "build",
    RunE: func(cmd *cobra.Command, args []string) error {
        return templates.BuildCache(afero.NewOsFs())
    },
}

var getCmd = &cobra.Command{
    Use: "get",
    RunE: func(cmd *cobra.Command, args []string) error {
        context, _ := cmd.Flags().GetString("context")
        step, _ := cmd.Flags().GetString("step")
        
        result, err := templates.GetStepTemplates(afero.NewOsFs(), context, step)
        if err != nil {
            return err
        }
        
        json.NewEncoder(os.Stdout).Encode(result)
        return nil
    },
}

var clearCmd = &cobra.Command{
    Use: "clear",
    RunE: func(cmd *cobra.Command, args []string) error {
        return afero.NewOsFs().Remove(".mission/libraries/cache.json")
    },
}

func init() {
    getCmd.Flags().String("context", "", "Context (plan/apply/complete)")
    getCmd.Flags().String("step", "", "Step (step1/step2/etc)")
    getCmd.MarkFlagRequired("context")
    getCmd.MarkFlagRequired("step")
    
    cacheCmd.AddCommand(buildCmd, getCmd, clearCmd)
    templatesCmd.AddCommand(cacheCmd)
    rootCmd.AddCommand(templatesCmd)
}
```

### 3. Updated Prompt Template

**File**: `internal/templates/prompts/m.plan.md`

**Step 1 Changes**:
```markdown
1. **Load Templates**: Run `m templates cache get --context plan --step step1`
2. **Analyze Intent**: Use templates["analysis/intent"]
3. **Verify Clarity**: Use templates["analysis/clarification"]
4. **If Ambiguous**: Use templates["missions/clarification"] and templates["displays/plan-clarification"]
```

**Step 2 Changes**:
```markdown
1. **Load Templates**: Run `m templates cache get --context plan --step step2`
2. **Duplication Check**: Use templates["analysis/duplication"]
3. **Domain Identification**: Use templates["analysis/domain"]
```

**Step 4 Changes**:
```markdown
**If Track 1**: Run `m templates cache get --context plan --step step4_track1`
**If Track 4**: Run `m templates cache get --context plan --step step4_track4`
```

**Step 6 Changes**:
```markdown
1. **Load Templates**: Run `m templates cache get --context plan --step step6`
2. **Output**: Use templates["displays/plan-success"]
```

---

## Benefits

1. **Single-Step Loading**: Only loads templates needed for current step
2. **Zero Waste**: 100% template utilization per step
3. **Multi-Command**: Same cache for plan/apply/complete
4. **Auto-Rebuild**: Missing cache triggers automatic rebuild
5. **Performance**: Single cache read vs multiple file reads

---

## Implementation Checklist

- [ ] Create `internal/templates/cache.go`
- [ ] Create `cmd/templates.go`
- [ ] Update `m.plan.md` to use cache commands
- [ ] Add cache build to `m init`
- [ ] Add `cache.json` to `.gitignore`
- [ ] Test all context.step combinations
- [ ] Update documentation

---

## Testing

```bash
# Build cache
m templates cache build

# Test each step
m templates cache get --context plan --step step1
m templates cache get --context plan --step step2
m templates cache get --context plan --step step4_track1
m templates cache get --context plan --step step6

# Verify JSON output structure
m templates cache get --context plan --step step1 | jq '.templates | keys'

# Clear and verify auto-rebuild
m templates cache clear
m templates cache get --context plan --step step1  # Should auto-rebuild
```
