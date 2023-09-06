package utilaio

const (
	EMPTY_STRING               = ""
	COMMA                      = ","
	COLON                      = ":"
	AT_SIGN                    = "@"
	NEW_LINE                   = "\n"
	EQUAL_TO                   = "="
	SPACE                      = " "
	FORWARD_SLASH_AS_STRING    = "/"
	FORWARD_SLASH_AS_CHARACTER = '/'
	BACKWARD_SLASH             = '\\'
	HYPHEN                     = "-"
	QUESTION_MARK              = "?"
	AMPERSAND                  = "&"
)

const (
	MINUTES_IN_1_DAY      = 1440
	SECONDS_IN_1_MINUTE   = 60
	SECONDS_IN_2_MINUTES  = 120
	SECONDS_IN_5_MINUTES  = 300
	SECONDS_IN_15_MINUTES = 900
	SECONDS_IN_1_HOUR     = 3600
)

const (
	JSON_FILE_EXTENSION          = ".json"
	WINDOWS_EXECUTABLE_EXTENSION = ".exe"
	WINDOWS_MSI_EXTENSION        = ".msi"
	ZIP_FILE_EXTENSION           = ".zip"
	CHECKSUM_FILE_EXTENSION      = ".sha256"
)

const (
	GOOS_WINDOWS = "windows"
	GOOS_LINUX   = "linux"
)

const (
	GOOS_386   = "386"
	GOOS_AMD64 = "amd64"
	GOOS_ARM   = "arm"
	GOOS_ARM64 = "arm64"
)

var (
	JSON_NULL *string
)
