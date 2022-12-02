package web

import (
	"bytes"
	"regexp"
	"strings"
)

// UserAgent struct containg all determined datra from parsed user-agent string
type UserAgent struct {
	Name       string `json:"name"`
	Version    string `json:"version"`
	OS         string `json:"os"`
	OSVersion  string `json:"os_version"`
	Device     string `json:"device,omitempty"`
	DeviceType string `json:"device_type,omitempty"`
	Mobile     bool   `json:"mobile,omitempty"`
	Tablet     bool   `json:"tablet,omitempty"`
	Desktop    bool   `json:"desktop,omitempty"`
	Script     bool   `json:"script,omitempty"`
	Bot        bool   `json:"bot,omitempty"`
	URL        string `json:"url,omitempty"`
	String     string `json:"raw,omitempty"`
}

var ignore = map[string]struct{}{
	"KHTML, like Gecko": struct{}{},
	"U":                 struct{}{},
	"compatible":        struct{}{},
	"Mozilla":           struct{}{},
	"WOW64":             struct{}{},
}

// Constants for browsers and operating systems for easier comparation
const (
	Windows      = "Windows"
	WindowsPhone = "Windows Phone"
	Android      = "Android"
	MacOS        = "macOS"
	IOS          = "iOS"
	Linux        = "Linux"

	Opera            = "Opera"
	OperaMini        = "Opera Mini"
	OperaTouch       = "Opera Touch"
	Chrome           = "Chrome"
	Firefox          = "Firefox"
	InternetExplorer = "Internet Explorer"
	Safari           = "Safari"
	Edge             = "Edge"
	Vivaldi          = "Vivaldi"

	Googlebot           = "Googlebot"
	Twitterbot          = "Twitterbot"
	FacebookExternalHit = "facebookexternalhit"
	Applebot            = "Applebot"
)

// Parse user agent string returning UserAgent struct
func Parse(userAgent string) UserAgent {
	ua := UserAgent{
		String: userAgent,
	}

	tokens := parse(userAgent)
	isTablet := false
	isDesktop := false
	isMobile := false
	isBot := false
	// check is there URL
	for k := range tokens {
		if strings.HasPrefix(k, "http://") || strings.HasPrefix(k, "https://") {
			ua.URL = k
			delete(tokens, k)
			break
		}
	}

	// OS lookup
	switch {
	case tokens.exists("Android"):
		ua.OS = Android
		ua.OSVersion = tokens[Android]
		for s := range tokens {
			if strings.HasSuffix(s, "Build") {
				ua.Device = strings.TrimSpace(s[:len(s)-5])
				isTablet = strings.Contains(strings.ToLower(ua.Device), "tablet")
			}
		}

	case tokens.exists("iPhone"):
		ua.OS = IOS
		ua.OSVersion = tokens.findMacOSVersion()
		ua.Device = "iPhone"
		isMobile = true

	case tokens.exists("iPad"):
		ua.OS = IOS
		ua.OSVersion = tokens.findMacOSVersion()
		ua.Device = "iPad"
		isTablet = true

	case tokens.exists("Windows NT"):
		ua.OS = Windows
		ua.OSVersion = tokens["Windows NT"]
		isDesktop = true

	case tokens.exists("Windows Phone OS"):
		ua.OS = WindowsPhone
		ua.OSVersion = tokens["Windows Phone OS"]
		isMobile = true

	case tokens.exists("Macintosh"):
		ua.OS = MacOS
		ua.OSVersion = tokens.findMacOSVersion()
		isDesktop = true

	case tokens.exists("Linux"):
		ua.OS = Linux
		ua.OSVersion = tokens[Linux]
		isDesktop = true

	}

	// for s, val := range sys {
	// 	fmt.Println(s, "--", val)
	// }

	switch {

	case tokens.exists("Googlebot"):
		ua.Name = Googlebot
		ua.Version = tokens[Googlebot]
		ua.Bot = true
		isMobile = tokens.existsAny("Mobile", "Mobile Safari")

	case tokens.exists("Applebot"):
		ua.Name = Applebot
		ua.Version = tokens[Applebot]
		ua.Bot = true
		isMobile = tokens.existsAny("Mobile", "Mobile Safari")
		ua.OS = ""

	case tokens["Opera Mini"] != "":
		ua.Name = OperaMini
		ua.Version = tokens[OperaMini]
		isMobile = true

	case tokens["OPR"] != "":
		ua.Name = Opera
		ua.Version = tokens["OPR"]
		isMobile = tokens.existsAny("Mobile", "Mobile Safari")

	case tokens["OPT"] != "":
		ua.Name = OperaTouch
		ua.Version = tokens["OPT"]
		isMobile = tokens.existsAny("Mobile", "Mobile Safari")

	// Opera on iOS
	case tokens["OPiOS"] != "":
		ua.Name = Opera
		ua.Version = tokens["OPiOS"]
		isMobile = tokens.existsAny("Mobile", "Mobile Safari")

	// Chrome on iOS
	case tokens["CriOS"] != "":
		ua.Name = Chrome
		ua.Version = tokens["CriOS"]
		isMobile = tokens.existsAny("Mobile", "Mobile Safari")

	// Firefox on iOS
	case tokens["FxiOS"] != "":
		ua.Name = Firefox
		ua.Version = tokens["FxiOS"]
		isMobile = tokens.existsAny("Mobile", "Mobile Safari")

	case tokens["Firefox"] != "":
		ua.Name = Firefox
		ua.Version = tokens[Firefox]
		_, isMobile = tokens["Mobile"]
		_, ua.Tablet = tokens["Tablet"]

	case tokens["Vivaldi"] != "":
		ua.Name = Vivaldi
		ua.Version = tokens[Vivaldi]

	case tokens.exists("MSIE"):
		ua.Name = InternetExplorer
		ua.Version = tokens["MSIE"]

	case tokens["EdgiOS"] != "":
		ua.Name = Edge
		ua.Version = tokens["EdgiOS"]
		isMobile = tokens.existsAny("Mobile", "Mobile Safari")

	case tokens["Edge"] != "":
		ua.Name = Edge
		ua.Version = tokens["Edge"]
		isMobile = tokens.existsAny("Mobile", "Mobile Safari")

	case tokens["Edg"] != "":
		ua.Name = Edge
		ua.Version = tokens["Edg"]
		isMobile = tokens.existsAny("Mobile", "Mobile Safari")

	case tokens["EdgA"] != "":
		ua.Name = Edge
		ua.Version = tokens["EdgA"]
		isMobile = tokens.existsAny("Mobile", "Mobile Safari")

	case tokens["bingbot"] != "":
		ua.Name = "Bingbot"
		ua.Version = tokens["bingbot"]
		isMobile = tokens.existsAny("Mobile", "Mobile Safari")

	case tokens["SamsungBrowser"] != "":
		ua.Name = "Samsung Browser"
		ua.Version = tokens["SamsungBrowser"]
		isMobile = tokens.existsAny("Mobile", "Mobile Safari")

	// if chrome and Safari defined, find any other tokensent descr
	case tokens.exists(Chrome) && tokens.exists(Safari):
		name := tokens.findBestMatch(true)
		if name != "" {
			ua.Name = name
			ua.Version = tokens[name]
			break
		}
		fallthrough

	case tokens.exists("Chrome"):
		ua.Name = Chrome
		ua.Version = tokens["Chrome"]
		isMobile = tokens.existsAny("Mobile", "Mobile Safari")

	case tokens.exists("Safari"):
		ua.Name = Safari
		if v, ok := tokens["Version"]; ok {
			ua.Version = v
		} else {
			ua.Version = tokens["Safari"]
		}
		isMobile = tokens.existsAny("Mobile", "Mobile Safari")

	default:
		if ua.OS == "Android" && tokens["Version"] != "" {
			ua.Name = "Android browser"
			ua.Version = tokens["Version"]
			isMobile = true
		} else {
			if name := tokens.findBestMatch(false); name != "" {
				ua.Name = name
				ua.Version = tokens[name]
			} else {
				ua.Name = ua.String
			}
			ua.Bot = strings.Contains(strings.ToLower(ua.Name), "bot")
			isMobile = tokens.existsAny("Mobile", "Mobile Safari")
		}
	}

	// if tablet, switch mobile to off
	if ua.Tablet {
		isMobile = false
	}

	// if not already bot, check some popular bots and weather URL is set
	botDetector := NewBotDetector()
	isBot = botDetector.IsBot(userAgent)

	if !isBot {
		switch ua.Name {
		case Twitterbot, FacebookExternalHit:
			isBot = true
		}
	}
	if isMobile {
		ua.DeviceType = "Mobile"
	}
	if isTablet {
		ua.DeviceType = "Tablet"
	}
	if isDesktop {
		ua.DeviceType = "Desktop"
	}
	if isBot {
		ua.DeviceType = "Bot"
	}
	isScript := !(isMobile || isTablet || isDesktop || isBot)
	if isScript {
		ua.DeviceType = "Script"
	}
	return ua
}

func parse(userAgent string) (clients properties) {
	clients = make(map[string]string)
	slash := false
	isURL := false
	var buff, val bytes.Buffer
	addToken := func() {
		if buff.Len() != 0 {
			s := strings.TrimSpace(buff.String())
			if _, ign := ignore[s]; !ign {
				if isURL {
					s = strings.TrimPrefix(s, "+")
				}

				if val.Len() == 0 { // only if value don't exists
					var ver string
					s, ver = checkVer(s) // determin version string and split
					clients[s] = ver
				} else {
					clients[s] = strings.TrimSpace(val.String())
				}
			}
		}
		buff.Reset()
		val.Reset()
		slash = false
		isURL = false
	}

	parOpen := false

	bua := []byte(userAgent)
	for i, c := range bua {

		//fmt.Println(string(c), c)
		switch {
		case c == 41: // )
			addToken()
			parOpen = false

		case parOpen && c == 59: // ;
			addToken()

		case c == 40: // (
			addToken()
			parOpen = true

		case slash && c == 32:
			addToken()

		case slash:
			val.WriteByte(c)

		case c == 47 && !isURL: //   /
			if i != len(bua)-1 && bua[i+1] == 47 && (bytes.HasSuffix(buff.Bytes(), []byte("http:")) || bytes.HasSuffix(buff.Bytes(), []byte("https:"))) {
				buff.WriteByte(c)
				isURL = true
			} else {
				slash = true
			}

		default:
			buff.WriteByte(c)
		}
	}
	addToken()

	return clients
}

func checkVer(s string) (name, v string) {
	i := strings.LastIndex(s, " ")
	if i == -1 {
		return s, ""
	}

	//v = s[i+1:]

	switch s[:i] {
	case "Linux", "Windows NT", "Windows Phone OS", "MSIE", "Android":
		return s[:i], s[i+1:]
	default:
		return s, ""
	}

	// for _, c := range v {
	// 	if (c >= 48 && c <= 57) || c == 46 {
	// 	} else {
	// 		return s, ""
	// 	}
	// }

	// return s[:i], s[i+1:]

}

type properties map[string]string

func (p properties) exists(key string) bool {
	_, ok := p[key]
	return ok
}

func (p properties) existsAny(keys ...string) bool {
	for _, k := range keys {
		if _, ok := p[k]; ok {
			return true
		}
	}
	return false
}

func (p properties) findMacOSVersion() string {
	for k, v := range p {
		if strings.Contains(k, "OS") {
			if ver := findVersion(v); ver != "" {
				return ver
			} else if ver = findVersion(k); ver != "" {
				return ver
			}
		}
	}
	return ""
}

// findBestMatch from the rest of the bunch
// in first cycle only return key vith version value
// if withVerValue is false, do another cycle and return any token
func (p properties) findBestMatch(withVerOnly bool) string {
	n := 2
	if withVerOnly {
		n = 1
	}
	for i := 0; i < n; i++ {
		for k, v := range p {
			switch k {
			case Chrome, Firefox, Safari, "Version", "Mobile", "Mobile Safari", "Mozilla", "AppleWebKit", "Windows NT", "Windows Phone OS", Android, "Macintosh", Linux, "GSA":
			default:
				if i == 0 {
					if v != "" { // in first check, only return  keys with value
						return k
					}
				} else {
					return k
				}
			}
		}
	}
	return ""
}

var rxMacOSVer = regexp.MustCompile(`[_\\d\\.]+`)

func findVersion(s string) string {
	if ver := rxMacOSVer.FindString(s); ver != "" {
		return strings.Replace(ver, "_", ".", -1)
	}
	return ""
}

// IsWindows shorthand function to check if OS == Windows
func (ua UserAgent) IsWindows() bool {
	return ua.OS == Windows
}

// IsAndroid shorthand function to check if OS == Android
func (ua UserAgent) IsAndroid() bool {
	return ua.OS == Android
}

// IsMacOS shorthand function to check if OS == MacOS
func (ua UserAgent) IsMacOS() bool {
	return ua.OS == MacOS
}

// IsIOS shorthand function to check if OS == IOS
func (ua UserAgent) IsIOS() bool {
	return ua.OS == IOS
}

// IsLinux shorthand function to check if OS == Linux
func (ua UserAgent) IsLinux() bool {
	return ua.OS == Linux
}

// IsOpera shorthand function to check if Name == Opera
func (ua UserAgent) IsOpera() bool {
	return ua.Name == Opera
}

// IsOperaMini shorthand function to check if Name == Opera Mini
func (ua UserAgent) IsOperaMini() bool {
	return ua.Name == OperaMini
}

// IsChrome shorthand function to check if Name == Chrome
func (ua UserAgent) IsChrome() bool {
	return ua.Name == Chrome
}

// IsFirefox shorthand function to check if Name == Firefox
func (ua UserAgent) IsFirefox() bool {
	return ua.Name == Firefox
}

// IsInternetExplorer shorthand function to check if Name == Internet Explorer
func (ua UserAgent) IsInternetExplorer() bool {
	return ua.Name == InternetExplorer
}

// IsSafari shorthand function to check if Name == Safari
func (ua UserAgent) IsSafari() bool {
	return ua.Name == Safari
}

// IsEdge shorthand function to check if Name == Edge
func (ua UserAgent) IsEdge() bool {
	return ua.Name == Edge
}

// IsGooglebot shorthand function to check if Name == Googlebot
func (ua UserAgent) IsGooglebot() bool {
	return ua.Name == Googlebot
}

// IsTwitterbot shorthand function to check if Name == Twitterbot
func (ua UserAgent) IsTwitterbot() bool {
	return ua.Name == Twitterbot
}

// IsFacebookbot shorthand function to check if Name == FacebookExternalHit
func (ua UserAgent) IsFacebookbot() bool {
	return ua.Name == FacebookExternalHit
}
