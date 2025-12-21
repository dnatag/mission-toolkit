# MOCK FACTORY

## PURPOSE
Centralized factory for creating consistent AI response and filesystem mocks across all unit tests.

## MOCK FACTORY ABSTRACTION

```javascript
export class MockFactory {
  constructor() {
    this.presets = new Map()
    this.setupPresets()
  }

  // Setup common mock presets
  setupPresets() {
    // AI Response Presets
    this.presets.set('simple-feature', {
      type: 'aiResponse',
      data: {
        intent: 'Add user profile endpoint with CRUD operations',
        suggestedFiles: ['routes/profile.js', 'controllers/profile.js'],
        domainMultipliers: [],
        complexity: 'standard'
      }
    })

    this.presets.set('security-feature', {
      type: 'aiResponse',
      data: {
        intent: 'Add JWT authentication with role-based access control',
        suggestedFiles: ['auth/jwt.js', 'auth/rbac.js', 'middleware/auth.js'],
        domainMultipliers: ['high-risk-integration', 'regulatory-security'],
        complexity: 'security-sensitive'
      }
    })

    this.presets.set('performance-feature', {
      type: 'aiResponse',
      data: {
        intent: 'Optimize search with sub-second response and caching',
        suggestedFiles: ['search/engine.js', 'cache/redis.js', 'routes/search.js'],
        domainMultipliers: ['performance-critical'],
        complexity: 'performance-critical'
      }
    })

    // Filesystem Presets
    this.presets.set('clean-project', {
      type: 'filesystem',
      data: {
        files: new Map([
          ['.mission/governance.md', 'governance content'],
          ['.mission/metrics.md', 'metrics content'],
          ['.mission/backlog.md', 'backlog content']
        ]),
        missionExists: false
      }
    })

    this.presets.set('active-mission', {
      type: 'filesystem',
      data: {
        files: new Map([
          ['.mission/mission.md', `# MISSION\ntype: WET\ntrack: 2\nstatus: active`]
        ]),
        missionExists: true
      }
    })
  }

  // Create AI response mock
  createAIResponseMock(preset, overrides = {}) {
    const baseData = this.presets.get(preset)?.data || {}
    
    return {
      ...baseData,
      ...overrides,
      // Add mock metadata
      _mockType: 'aiResponse',
      _preset: preset,
      _timestamp: Date.now()
    }
  }

  // Create filesystem mock
  createFilesystemMock(preset, overrides = {}) {
    const baseData = this.presets.get(preset)?.data || { files: new Map() }
    
    class MockFileSystem {
      constructor(initialData) {
        this.files = new Map(initialData.files)
        this.capturedWrites = new Map()
        this.capturedReads = []
        this.existsChecks = []
      }

      exists(path) {
        this.existsChecks.push(path)
        return this.files.has(path)
      }

      readFile(path) {
        this.capturedReads.push(path)
        if (!this.files.has(path)) {
          throw new Error(`File not found: ${path}`)
        }
        return this.files.get(path)
      }

      writeFile(path, content) {
        this.capturedWrites.set(path, content)
        this.files.set(path, content)
        return true
      }

      getCapturedWrites() {
        return Object.fromEntries(this.capturedWrites)
      }

      getReadAttempts() {
        return this.capturedReads
      }

      getExistenceChecks() {
        return this.existsChecks
      }

      reset() {
        this.capturedWrites.clear()
        this.capturedReads.length = 0
        this.existsChecks.length = 0
      }
    }

    const mockFS = new MockFileSystem(baseData)
    
    // Apply overrides
    if (overrides.additionalFiles) {
      for (const [path, content] of Object.entries(overrides.additionalFiles)) {
        mockFS.files.set(path, content)
      }
    }

    return mockFS
  }

  // Create execution logger mock
  createLoggerMock(config = {}) {
    const logs = []
    
    return {
      logPromptInvocation: (prompt, input) => {
        logs.push({ type: 'PROMPT_INVOCATION', prompt, input, timestamp: Date.now() })
      },
      
      logMockAIResponse: (responseType, data) => {
        logs.push({ type: 'MOCK_AI_RESPONSE', responseType, data, timestamp: Date.now() })
      },
      
      logPromptLogicExecution: (operation, input, output) => {
        logs.push({ type: 'PROMPT_LOGIC_EXECUTION', operation, input, output, timestamp: Date.now() })
      },
      
      logDecisionPoint: (decision, reasoning, result) => {
        logs.push({ type: 'DECISION_POINT', decision, reasoning, result, timestamp: Date.now() })
      },
      
      getLogs: () => [...logs],
      
      getExecutionTrace: () => logs.map((log, index) => 
        `[${index + 1}] ${log.type}: ${log.prompt || log.operation || log.decision}`
      ).join('\n'),
      
      clear: () => logs.length = 0
    }
  }

  // Create command runner mock
  createCommandMock(responses = {}) {
    const defaultResponses = {
      'npm test': { exitCode: 0, output: 'All tests passed' },
      'go test ./...': { exitCode: 0, output: 'PASS' },
      'grep -r pattern': { exitCode: 0, output: 'pattern found' }
    }

    return {
      execute: (command) => {
        const response = responses[command] || defaultResponses[command] || {
          exitCode: 0,
          output: `Mock output for: ${command}`
        }
        
        return {
          ...response,
          command,
          timestamp: Date.now()
        }
      }
    }
  }

  // Create complete test environment
  createTestEnvironment(scenario) {
    const aiMock = this.createAIResponseMock(scenario.aiPreset, scenario.aiOverrides)
    const fsMock = this.createFilesystemMock(scenario.fsPreset, scenario.fsOverrides)
    const loggerMock = this.createLoggerMock(scenario.loggerConfig)
    const commandMock = this.createCommandMock(scenario.commandResponses)

    return {
      aiResponse: aiMock,
      fileSystem: fsMock,
      logger: loggerMock,
      commandRunner: commandMock,
      
      // Convenience methods
      reset: () => {
        fsMock.reset()
        loggerMock.clear()
      },
      
      getCaptures: () => ({
        filesWritten: fsMock.getCapturedWrites(),
        filesRead: fsMock.getReadAttempts(),
        logs: loggerMock.getLogs()
      })
    }
  }

  // Register custom preset
  registerPreset(name, type, data) {
    this.presets.set(name, { type, data })
  }

  // Get available presets
  getAvailablePresets(type = null) {
    const presets = Array.from(this.presets.entries())
    
    if (type) {
      return presets
        .filter(([_, preset]) => preset.type === type)
        .map(([name, _]) => name)
    }
    
    return presets.map(([name, preset]) => ({ name, type: preset.type }))
  }
}

// Singleton instance
export const mockFactory = new MockFactory()
```

## USAGE EXAMPLES

```javascript
import { mockFactory } from './mock-factory.md'

// Create simple test environment
const env = mockFactory.createTestEnvironment({
  aiPreset: 'simple-feature',
  fsPreset: 'clean-project'
})

// Create custom AI response
const customAI = mockFactory.createAIResponseMock('security-feature', {
  intent: 'Custom security implementation',
  suggestedFiles: ['custom/security.js']
})

// Create filesystem with additional files
const customFS = mockFactory.createFilesystemMock('clean-project', {
  additionalFiles: {
    'custom/file.js': 'custom content'
  }
})

// Register custom preset
mockFactory.registerPreset('payment-feature', 'aiResponse', {
  intent: 'Add payment processing',
  suggestedFiles: ['payments/stripe.js', 'payments/paypal.js'],
  domainMultipliers: ['high-risk-integration', 'regulatory-security']
})
```
