# TEST RUNNER

## PURPOSE
Unified test execution framework for AI prompt unit tests with standardized lifecycle and reporting.

## TEST RUNNER ABSTRACTION

```javascript
export class TestRunner {
  constructor(config = {}) {
    this.config = {
      verbose: config.verbose || false,
      stopOnFailure: config.stopOnFailure || false,
      timeout: config.timeout || 5000,
      ...config
    }
    this.results = {
      total: 0,
      passed: 0,
      failed: 0,
      skipped: 0,
      errors: []
    }
  }

  // Execute single test case
  async runTest(testCase) {
    const startTime = Date.now()
    
    try {
      // Setup phase
      if (testCase.setup) {
        await testCase.setup()
      }
      
      // Execution phase
      const result = await testCase.execute()
      
      // Assertion phase
      if (testCase.assertions) {
        for (const assertion of testCase.assertions) {
          assertion.validate(result)
        }
      }
      
      // Cleanup phase
      if (testCase.cleanup) {
        await testCase.cleanup()
      }
      
      this.results.passed++
      return {
        name: testCase.name,
        status: 'PASSED',
        duration: Date.now() - startTime,
        result
      }
      
    } catch (error) {
      this.results.failed++
      this.results.errors.push({
        test: testCase.name,
        error: error.message,
        stack: error.stack
      })
      
      return {
        name: testCase.name,
        status: 'FAILED',
        duration: Date.now() - startTime,
        error: error.message
      }
    }
  }

  // Execute test suite
  async runSuite(testSuite) {
    const suiteResults = {
      name: testSuite.name,
      tests: [],
      summary: { passed: 0, failed: 0, skipped: 0 }
    }
    
    for (const testCase of testSuite.tests) {
      this.results.total++
      
      if (testCase.skip) {
        this.results.skipped++
        suiteResults.tests.push({
          name: testCase.name,
          status: 'SKIPPED'
        })
        continue
      }
      
      const result = await this.runTest(testCase)
      suiteResults.tests.push(result)
      
      if (result.status === 'PASSED') {
        suiteResults.summary.passed++
      } else {
        suiteResults.summary.failed++
        
        if (this.config.stopOnFailure) {
          break
        }
      }
    }
    
    return suiteResults
  }

  // Execute multiple test suites
  async runAll(testSuites) {
    const allResults = []
    
    for (const suite of testSuites) {
      const suiteResult = await this.runSuite(suite)
      allResults.push(suiteResult)
    }
    
    return {
      suites: allResults,
      summary: this.results,
      success: this.results.failed === 0
    }
  }

  // Generate test report
  generateReport(results) {
    const report = {
      timestamp: new Date().toISOString(),
      summary: results.summary,
      success: results.success,
      details: results.suites.map(suite => ({
        suite: suite.name,
        passed: suite.summary.passed,
        failed: suite.summary.failed,
        skipped: suite.summary.skipped,
        failedTests: suite.tests
          .filter(test => test.status === 'FAILED')
          .map(test => ({ name: test.name, error: test.error }))
      }))
    }
    
    return report
  }
}

// Test case builder
export class TestCaseBuilder {
  constructor(name) {
    this.testCase = { name }
  }
  
  setup(setupFn) {
    this.testCase.setup = setupFn
    return this
  }
  
  execute(executeFn) {
    this.testCase.execute = executeFn
    return this
  }
  
  assert(assertions) {
    this.testCase.assertions = Array.isArray(assertions) ? assertions : [assertions]
    return this
  }
  
  cleanup(cleanupFn) {
    this.testCase.cleanup = cleanupFn
    return this
  }
  
  skip(reason) {
    this.testCase.skip = true
    this.testCase.skipReason = reason
    return this
  }
  
  build() {
    return this.testCase
  }
}

// Test suite builder
export class TestSuiteBuilder {
  constructor(name) {
    this.suite = { name, tests: [] }
  }
  
  addTest(testCase) {
    this.suite.tests.push(testCase)
    return this
  }
  
  addTests(testCases) {
    this.suite.tests.push(...testCases)
    return this
  }
  
  build() {
    return this.suite
  }
}
```

## USAGE EXAMPLE

```javascript
import { TestRunner, TestCaseBuilder, TestSuiteBuilder } from './test-runner.md'

// Create test cases
const testCase1 = new TestCaseBuilder('Parse user intent')
  .setup(() => ({ input: 'add user auth' }))
  .execute(({ input }) => promptLogic.parseUserIntent(input))
  .assert([
    result => assert(result.wordCount === 3),
    result => assert(result.isVague === false)
  ])
  .build()

// Create test suite
const suite = new TestSuiteBuilder('Prompt Logic Tests')
  .addTest(testCase1)
  .build()

// Run tests
const runner = new TestRunner({ verbose: true })
const results = await runner.runAll([suite])
const report = runner.generateReport(results)
```
