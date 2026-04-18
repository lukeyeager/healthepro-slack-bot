package scheduler

// The scheduler's logic is now a thin loop around the healthepro client and
// store — both of which have their own test suites.  Week boundary logic lives
// in internal/week and is tested there.  Nothing remains to unit-test here.
