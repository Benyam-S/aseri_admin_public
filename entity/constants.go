package entity

// AppSecretKeyName is a constant that defines the app secret key name
const AppSecretKeyName = "aseri_secret_key"

// AppCookieName is a constant that defines the cookie name
const AppCookieName = "aseri_cookie_name"

// RoleStaff is a constant that defines a staff role for a staff member table
const RoleStaff = "Staff"

// RoleAdmin is a constant that defines a admin role for a staff member table
const RoleAdmin = "Admin"

// RoleAny is a constant that defines role to be any type
const RoleAny = "Any"

// UserCategoryAgent is a constant that holds the agent user category
const UserCategoryAgent = "Agent"

// UserCategoryAseri is a constant that holds the aseri user category
const UserCategoryAseri = "Aseri"

// UserCategoryAny is a constant that defines a user category to be of any type
const UserCategoryAny = "Any"

// UserCategoryJobSeeker is a constant that holds the job seeker user category
const UserCategoryJobSeeker = "JobSeeker"

// PostCategoryInternal is a constant that states the job is posted by internal staff member
const PostCategoryInternal = "Internal"

// PostCategoryUser is a constant that states the job is posted by aseri user
const PostCategoryUser = "User"

// PostCategoryExternal is a constant that states the job is posted by or from external third party
const PostCategoryExternal = "External"

// JobStatusPending is a constant that states a job is in pending state for approval
const JobStatusPending = "P"

// JobStatusOpened is a constant that states a job has been approved and open for application
const JobStatusOpened = "O"

// JobStatusClosed is a constant that states a job is in closed state
const JobStatusClosed = "C"

// JobStatusDecelined is a constant that states a job has been decelined
const JobStatusDecelined = "D"

// JobStatusAny is a constant that defines a job status to be of any type
const JobStatusAny = "Any"

// FeedbackSeen is a constant that states a feedback has been seen
const FeedbackSeen = "Seen"

// FeedbackUnseen is a constant that states a feedback hasn't been seen
const FeedbackUnseen = "Unseen"

// StartPush is a constant that states start push
const StartPush = "Start"

// ServerLogFile is a constant that holds the server log file name
const ServerLogFile = "server.log"

// BotLogFile is a constant that holds the bot log file name
const BotLogFile = "bot.log"

// PushForApproval is a constant that states push for approval key
const PushForApproval = "Approval"

// PushToChannel is a constant that states push to channel key
const PushToChannel = "Channel"

// PushToSubscribers is a constant that states push to subscribers key
const PushToSubscribers = "Subscribers"

// ValidWorkExperiences is a value list that holds all the valid work experience
var ValidWorkExperiences = []string{"0 year", "1 year", "2 years", "3 years", "4 years",
	"5 years", "6 years", "7 years", "8 years", "9 years", "10+ years"}

// ValidContactTypes is a value list that holds all the valid contact types
var ValidContactTypes = []string{"Via Telegram Account", "Send CV"}
