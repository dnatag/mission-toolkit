# M.CLARIFY TEST FIXTURE

## TEST_MODE
dry_run: true
validate_only: true

## SETUP_MISSION
```markdown
# MISSION

type: CLARIFICATION
track: TBD
iteration: 1
status: clarifying

## INTENT
Build a payment processing system

## NEED_CLARIFICATION
- [ ] Which payment providers should be supported (Stripe, PayPal, Square)?
- [ ] What security compliance is required (PCI DSS level)?
- [ ] Should it handle recurring payments or one-time only?

## PROVISIONAL_SCOPE
payments/processor.js
payments/providers/
config/payment-config.js
```

## INPUT
```text
We need Stripe and PayPal support, PCI DSS Level 1 compliance, and both one-time and recurring payments.
```

## EXPECTED_CLARIFIED_MISSION
```markdown
# MISSION

type: WET
track: 3
iteration: 1
status: planned

## INTENT
Build payment processing system with Stripe and PayPal support, PCI DSS Level 1 compliance, handling both one-time and recurring payments

## SCOPE
payments/stripe-processor.js
payments/paypal-processor.js
payments/recurring-manager.js
config/payment-config.js
middleware/pci-compliance.js

## PLAN
- [ ] Implement Stripe payment processor with PCI compliance
- [ ] Implement PayPal payment processor with PCI compliance
- [ ] Create recurring payment management system
- [ ] Add PCI DSS Level 1 compliance middleware
- [ ] Note: Allow duplication for initial implementation

## VERIFICATION
npm test -- --grep "payment"
```

## STEP_ASSERTIONS
- [ ] Mission type changed from CLARIFICATION to WET
- [ ] Track reassessed to 3 (due to security compliance)
- [ ] All clarification questions addressed in INTENT
- [ ] SCOPE expanded based on clarifications
- [ ] Status changed to "planned"

## VALIDATION_INSTRUCTIONS
```
Execute @m.clarify in validation mode:
- Process clarification responses
- Update mission structure accordingly
- DO NOT write .mission/mission.md to disk
- Return validation results only
```
