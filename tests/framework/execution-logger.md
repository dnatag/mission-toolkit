# EXECUTION LOGGER

## PURPOSE
Reusable execution logging framework for capturing detailed test execution traces with structured output and analysis.

## EXECUTION LOGGER ABSTRACTION

```javascript
export class ExecutionLogger {
  constructor(config = {}) {
    this.config = {
      verbose: config.verbose || false,
      maxLogs: config.maxLogs || 1000,
      timestampFormat: config.timestampFormat || 'ISO',
      ...config
    }
    this.logs = []
    this.stepCounter = 0
    this.startTime = Date.now()
  }

  // Core logging method
  log(type, data) {
    this.stepCounter++
    
    const logEntry = {
      step: this.stepCounter,
      type,
      timestamp: this.formatTimestamp(Date.now()),
      duration: Date.now() - this.startTime,
      ...data
    }
    
    this.logs.push(logEntry)
    
    // Enforce max logs limit
    if (this.logs.length > this.config.maxLogs) {
      this.logs.shift()
    }
    
    if (this.config.verbose) {
      this.printLog(logEntry)
    }
    
    return logEntry
  }

  // Specific log types
  logPromptInvocation(prompt, input) {
    return this.log('PROMPT_INVOCATION', {
      prompt,
      input: this.truncateInput(input),
      inputLength: input?.length || 0
    })
  }

  logMockAIResponse(responseType, data) {
    return this.log('MOCK_AI_RESPONSE', {
      responseType,
      data: this.sanitizeData(data),
      dataSize: JSON.stringify(data).length
    })
  }

  logPromptLogicExecution(operation, input, output) {
    return this.log('PROMPT_LOGIC_EXECUTION', {
      operation,
      input: this.sanitizeData(input),
      output: this.sanitizeData(output),
      success: output !== null && output !== undefined
    })
  }

  logDecisionPoint(decision, reasoning, result) {
    return this.log('DECISION_POINT', {
      decision,
      reasoning,
      result,
      confidence: this.calculateConfidence(reasoning, result)
    })
  }

  logFileSystemOperation(operation, path, metadata = {}) {
    return this.log('FILESYSTEM_OPERATION', {
      operation,
      path,
      ...metadata,
      pathDepth: path.split('/').length
    })
  }

  logAssertion(assertion, passed, error = null) {
    return this.log('ASSERTION', {
      assertion,
      passed,
      error: error?.message || null,
      severity: passed ? 'INFO' : 'ERROR'
    })
  }

  logTestPhase(phase, testName, metadata = {}) {
    return this.log('TEST_PHASE', {
      phase, // 'SETUP', 'EXECUTE', 'ASSERT', 'CLEANUP'
      testName,
      ...metadata
    })
  }

  // Utility methods
  formatTimestamp(timestamp) {
    switch (this.config.timestampFormat) {
      case 'ISO':
        return new Date(timestamp).toISOString()
      case 'relative':
        return `+${timestamp - this.startTime}ms`
      default:
        return timestamp
    }
  }

  truncateInput(input, maxLength = 100) {
    if (!input || typeof input !== 'string') return input
    return input.length > maxLength ? input.substring(0, maxLength) + '...' : input
  }

  sanitizeData(data) {
    if (data === null || data === undefined) return data
    
    try {
      const serialized = JSON.stringify(data)
      return serialized.length > 500 ? 
        JSON.parse(serialized.substring(0, 500) + '"}') : 
        data
    } catch {
      return String(data).substring(0, 100)
    }
  }

  calculateConfidence(reasoning, result) {
    // Simple heuristic for decision confidence
    const reasoningLength = reasoning?.length || 0
    const hasResult = result !== null && result !== undefined
    
    if (reasoningLength > 50 && hasResult) return 'HIGH'
    if (reasoningLength > 20 && hasResult) return 'MEDIUM'
    return 'LOW'
  }

  printLog(logEntry) {
    const icons = {
      'PROMPT_INVOCATION': '🚀',
      'MOCK_AI_RESPONSE': '🤖',
      'PROMPT_LOGIC_EXECUTION': '⚙️',
      'DECISION_POINT': '🎯',
      'FILESYSTEM_OPERATION': '📁',
      'ASSERTION': '✅',
      'TEST_PHASE': '📋'
    }
    
    const icon = icons[logEntry.type] || '📝'
    console.log(`[${logEntry.step}] ${icon} ${logEntry.type}: ${this.formatLogMessage(logEntry)}`)
  }

  formatLogMessage(logEntry) {
    switch (logEntry.type) {
      case 'PROMPT_INVOCATION':
        return `${logEntry.prompt} (${logEntry.inputLength} chars)`
      case 'MOCK_AI_RESPONSE':
        return `${logEntry.responseType} (${logEntry.dataSize} bytes)`
      case 'PROMPT_LOGIC_EXECUTION':
        return `${logEntry.operation} → ${logEntry.success ? 'SUCCESS' : 'FAILED'}`
      case 'DECISION_POINT':
        return `${logEntry.decision} → ${logEntry.result} (${logEntry.confidence})`
      case 'FILESYSTEM_OPERATION':
        return `${logEntry.operation} ${logEntry.path}`
      case 'ASSERTION':
        return `${logEntry.assertion} → ${logEntry.passed ? 'PASS' : 'FAIL'}`
      case 'TEST_PHASE':
        return `${logEntry.phase} - ${logEntry.testName}`
      default:
        return JSON.stringify(logEntry).substring(0, 100)
    }
  }

  // Analysis methods
  getExecutionTrace() {
    return this.logs.map(log => 
      `[${log.step}] ${log.type}: ${this.formatLogMessage(log)}`
    ).join('\n')
  }

  getExecutionSummary() {
    const summary = {
      totalSteps: this.stepCounter,
      duration: Date.now() - this.startTime,
      logTypes: {},
      phases: [],
      errors: []
    }
    
    this.logs.forEach(log => {
      // Count log types
      summary.logTypes[log.type] = (summary.logTypes[log.type] || 0) + 1
      
      // Track test phases
      if (log.type === 'TEST_PHASE') {
        summary.phases.push({ phase: log.phase, testName: log.testName, step: log.step })
      }
      
      // Collect errors
      if (log.type === 'ASSERTION' && !log.passed) {
        summary.errors.push({ step: log.step, assertion: log.assertion, error: log.error })
      }
    })
    
    return summary
  }

  getLogsByType(type) {
    return this.logs.filter(log => log.type === type)
  }

  getLogsByTimeRange(startTime, endTime) {
    return this.logs.filter(log => {
      const logTime = new Date(log.timestamp).getTime()
      return logTime >= startTime && logTime <= endTime
    })
  }

  // Export/Import methods
  exportLogs(format = 'json') {
    switch (format) {
      case 'json':
        return JSON.stringify({
          config: this.config,
          summary: this.getExecutionSummary(),
          logs: this.logs
        }, null, 2)
      case 'csv':
        const headers = 'step,type,timestamp,duration,details'
        const rows = this.logs.map(log => 
          `${log.step},${log.type},${log.timestamp},${log.duration},"${this.formatLogMessage(log)}"`
        )
        return [headers, ...rows].join('\n')
      default:
        return this.getExecutionTrace()
    }
  }

  // Utility methods
  clear() {
    this.logs = []
    this.stepCounter = 0
    this.startTime = Date.now()
  }

  clone() {
    const cloned = new ExecutionLogger(this.config)
    cloned.logs = [...this.logs]
    cloned.stepCounter = this.stepCounter
    cloned.startTime = this.startTime
    return cloned
  }
}

// Factory function
export function createLogger(config = {}) {
  return new ExecutionLogger(config)
}

// Singleton instance for global logging
export const globalLogger = new ExecutionLogger({ verbose: false })
```

## USAGE EXAMPLES

```javascript
import { createLogger, globalLogger } from './execution-logger.md'

// Create test-specific logger
const logger = createLogger({ verbose: true, maxLogs: 500 })

// Log test execution
logger.logTestPhase('SETUP', 'M.PLAN unit test')
logger.logPromptInvocation('@m.plan', 'add user authentication')
logger.logMockAIResponse('intent_analysis', { intent: 'Add auth system' })
logger.logPromptLogicExecution('calculateTrackComplexity', { files: 3 }, { track: 2 })
logger.logDecisionPoint('track_assignment', 'Based on file count and multipliers', 'Track 2')
logger.logAssertion('Track should be 2', true)
logger.logTestPhase('CLEANUP', 'M.PLAN unit test')

// Get execution summary
const summary = logger.getExecutionSummary()
console.log(`Test completed in ${summary.duration}ms with ${summary.totalSteps} steps`)

// Export logs
const jsonLogs = logger.exportLogs('json')
const csvLogs = logger.exportLogs('csv')
```
