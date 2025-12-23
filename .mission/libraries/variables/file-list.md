FILE_LIST = Estimate files needed based on intent and project language:

**Pattern Examples by Intent:**
- API endpoint: routes/{{name}}, controllers/{{name}}, handlers/{{name}}
- UI component: components/{{Name}}, views/{{name}}, templates/{{name}}
- Database: models/{{name}}, schemas/{{name}}, migrations/{{timestamp}}_{{name}}
- Authentication: auth/{{method}}, middleware/auth, guards/{{name}}
- Configuration: config/{{name}}, settings/{{name}}, env/{{name}}
- Testing: tests/{{name}}, spec/{{name}}, __tests__/{{name}} (don't count toward complexity)

**Adapt to Project Language:**
- Use existing file extensions (.js, .py, .go, .rs, .java, etc.)
- Follow project directory structure
- Match existing naming conventions

Format: One file path per line, no bullets or numbers