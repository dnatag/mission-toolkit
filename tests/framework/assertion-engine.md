# ASSERTION ENGINE

## PURPOSE
Standardized assertion framework for validating AI prompt logic results with fluent API and detailed error reporting.

## ASSERTION ENGINE ABSTRACTION

```javascript
export class AssertionEngine {
  constructor() {
    this.assertions = []
    this.results = { passed: 0, failed: 0, errors: [] }
  }

  // Create assertion builder
  expect(actual) {
    return new AssertionBuilder(actual, this)
  }

  // Record assertion result
  recordResult(assertion, passed, error = null) {
    this.assertions.push({
      description: assertion.description,
      passed,
      error,
      timestamp: Date.now()
    })

    if (passed) {
      this.results.passed++
    } else {
      this.results.failed++
      this.results.errors.push({
        assertion: assertion.description,
        error: error?.message || 'Assertion failed'
      })
    }
  }

  // Get assertion results
  getResults() {
    return {
      ...this.results,
      total: this.results.passed + this.results.failed,
      success: this.results.failed === 0,
      assertions: this.assertions
    }
  }

  // Reset for new test
  reset() {
    this.assertions = []
    this.results = { passed: 0, failed: 0, errors: [] }
  }
}

export class AssertionBuilder {
  constructor(actual, engine) {
    this.actual = actual
    this.engine = engine
    this.description = ''
  }

  // Set assertion description
  describe(description) {
    this.description = description
    return this
  }

  // Value equality assertions
  toBe(expected) {
    try {
      const passed = this.actual === expected
      if (!passed) {
        throw new Error(`Expected ${this.actual} to be ${expected}`)
      }
      this.engine.recordResult(this, true)
    } catch (error) {
      this.engine.recordResult(this, false, error)
    }
    return this
  }

  toEqual(expected) {
    try {
      const passed = JSON.stringify(this.actual) === JSON.stringify(expected)
      if (!passed) {
        throw new Error(`Expected ${JSON.stringify(this.actual)} to equal ${JSON.stringify(expected)}`)
      }
      this.engine.recordResult(this, true)
    } catch (error) {
      this.engine.recordResult(this, false, error)
    }
    return this
  }

  // Type assertions
  toBeType(expectedType) {
    try {
      const actualType = typeof this.actual
      const passed = actualType === expectedType
      if (!passed) {
        throw new Error(`Expected type ${expectedType}, got ${actualType}`)
      }
      this.engine.recordResult(this, true)
    } catch (error) {
      this.engine.recordResult(this, false, error)
    }
    return this
  }

  // Array assertions
  toContain(expected) {
    try {
      const passed = Array.isArray(this.actual) && this.actual.includes(expected)
      if (!passed) {
        throw new Error(`Expected array to contain ${expected}`)
      }
      this.engine.recordResult(this, true)
    } catch (error) {
      this.engine.recordResult(this, false, error)
    }
    return this
  }

  toHaveLength(expectedLength) {
    try {
      const actualLength = this.actual?.length
      const passed = actualLength === expectedLength
      if (!passed) {
        throw new Error(`Expected length ${expectedLength}, got ${actualLength}`)
      }
      this.engine.recordResult(this, true)
    } catch (error) {
      this.engine.recordResult(this, false, error)
    }
    return this
  }

  // Object assertions
  toHaveProperty(property, expectedValue = undefined) {
    try {
      const hasProperty = this.actual && this.actual.hasOwnProperty(property)
      if (!hasProperty) {
        throw new Error(`Expected object to have property ${property}`)
      }
      
      if (expectedValue !== undefined) {
        const actualValue = this.actual[property]
        const valueMatches = actualValue === expectedValue
        if (!valueMatches) {
          throw new Error(`Expected property ${property} to be ${expectedValue}, got ${actualValue}`)
        }
      }
      
      this.engine.recordResult(this, true)
    } catch (error) {
      this.engine.recordResult(this, false, error)
    }
    return this
  }

  // Boolean assertions
  toBeTruthy() {
    try {
      const passed = !!this.actual
      if (!passed) {
        throw new Error(`Expected ${this.actual} to be truthy`)
      }
      this.engine.recordResult(this, true)
    } catch (error) {
      this.engine.recordResult(this, false, error)
    }
    return this
  }

  toBeFalsy() {
    try {
      const passed = !this.actual
      if (!passed) {
        throw new Error(`Expected ${this.actual} to be falsy`)
      }
      this.engine.recordResult(this, true)
    } catch (error) {
      this.engine.recordResult(this, false, error)
    }
    return this
  }

  // Error assertions
  toThrow(expectedError = null) {
    try {
      let threwError = false
      let actualError = null
      
      try {
        if (typeof this.actual === 'function') {
          this.actual()
        }
      } catch (error) {
        threwError = true
        actualError = error
      }
      
      if (!threwError) {
        throw new Error('Expected function to throw an error')
      }
      
      if (expectedError && !actualError.message.includes(expectedError)) {
        throw new Error(`Expected error to contain "${expectedError}", got "${actualError.message}"`)
      }
      
      this.engine.recordResult(this, true)
    } catch (error) {
      this.engine.recordResult(this, false, error)
    }
    return this
  }

  // Custom assertion
  toSatisfy(predicate, errorMessage = 'Custom assertion failed') {
    try {
      const passed = predicate(this.actual)
      if (!passed) {
        throw new Error(errorMessage)
      }
      this.engine.recordResult(this, true)
    } catch (error) {
      this.engine.recordResult(this, false, error)
    }
    return this
  }
}

// Specialized assertions for prompt logic
export class PromptAssertions extends AssertionEngine {
  
  // Mission structure assertions
  expectValidMission(mission) {
    const builder = this.expect(mission).describe('Valid mission structure')
    
    // Check required fields
    builder.toHaveProperty('type')
    builder.toHaveProperty('track')
    builder.toHaveProperty('status')
    builder.toHaveProperty('intent')
    builder.toHaveProperty('scope')
    
    return this
  }

  // Track complexity assertions
  expectTrackCalculation(fileCount, multipliers, expectedTrack) {
    const result = promptLogic.calculateTrackComplexity(fileCount, multipliers)
    
    this.expect(result.finalTrack)
      .describe(`Track calculation for ${fileCount} files + ${multipliers.length} multipliers`)
      .toBe(expectedTrack)
    
    return this
  }

  // Domain multiplier assertions
  expectMultiplierDetection(text, expectedMultipliers) {
    const result = promptLogic.detectDomainMultipliers(text)
    
    for (const expected of expectedMultipliers) {
      this.expect(result.detected)
        .describe(`Multiplier detection for "${expected}"`)
        .toContain(expected)
    }
    
    return this
  }

  // Intent parsing assertions
  expectIntentParsing(input, expectedProperties) {
    const result = promptLogic.parseUserIntent(input)
    
    for (const [property, expectedValue] of Object.entries(expectedProperties)) {
      this.expect(result)
        .describe(`Intent parsing - ${property}`)
        .toHaveProperty(property, expectedValue)
    }
    
    return this
  }
}

// Factory function for creating assertion engines
export function createAssertions(type = 'standard') {
  switch (type) {
    case 'prompt':
      return new PromptAssertions()
    default:
      return new AssertionEngine()
  }
}
```

## USAGE EXAMPLES

```javascript
import { createAssertions } from './assertion-engine.md'

// Standard assertions
const assert = createAssertions()

assert.expect(result.track).describe('Track should be 2').toBe(2)
assert.expect(result.scope).describe('Scope should have 3 files').toHaveLength(3)
assert.expect(result.type).describe('Type should be WET').toBe('WET')

// Prompt-specific assertions
const promptAssert = createAssertions('prompt')

promptAssert.expectTrackCalculation(5, ['security'], 3)
promptAssert.expectMultiplierDetection('add payment processing', ['high-risk-integration'])
promptAssert.expectIntentParsing('add user auth', { wordCount: 3, isVague: false })

// Get results
const results = assert.getResults()
console.log(`Passed: ${results.passed}, Failed: ${results.failed}`)
```
