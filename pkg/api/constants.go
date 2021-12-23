package api

// CreateAPIAction represents the create api name
const CreateAPIAction = "create"

// JoinAPIAction represents the create api name
const JoinAPIAction = "join"

// SuccessReturnCode represents the api success return code
const SuccessReturnCode = "SUCCESS"

// FailedReturnCode represents the api ailed return code
const FailedReturnCode = "FAILED"

// ValidationErrorMessageKey represents the api validation error message key
const ValidationErrorMessageKey = "validationError"

// DuplicationWarningMessageKey represents the api duplication warning message key
const DuplicationWarningMessageKey = "duplicationWarning"

// NotFoundMessageKey represents the api not found message key
const NotFoundMessageKey = "notFound"

// EmptyMeetingNameMessage represents the api empty meeting name message
const EmptyMeetingNameMessage = "You must provide a meeting name"

// EmptyMeetingIDMessage represents the api empty meeting id message
const EmptyMeetingIDMessage = "You must provide a meeting ID"

// DuplicationWarningMessage represents the api duplication warning message
const DuplicationWarningMessage = "This conference was already in existence and may currently be in progress."

// NotFoundMeetingIDMessage represents the api not found message
const NotFoundMeetingIDMessage = "A meeting with that ID does not exist"
