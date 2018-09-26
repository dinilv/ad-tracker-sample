package v1

import (
	"net/url"
	"regexp"
	"strings"
)

var reg, regSpecial, RegNumbers, regEscapeString *regexp.Regexp

func init() {
	reg = regexp.MustCompile(".*{+.*}*.*")
	regSpecial = regexp.MustCompile("({\\w+})+")
	RegNumbers = regexp.MustCompile("^\\d+")
	regEscapeString = regexp.MustCompile("%7B.*%7D")
}

func ReplaceURLParameters(params map[string]string, templateUrl string) string {

	template, _ := url.Parse(templateUrl)
	scheme := template.Scheme
	if scheme == "" {
		scheme = "http"
		template, _ = url.Parse(scheme + "://" + templateUrl)
	}
	escapePath := template.EscapedPath()
	finalTemplate := scheme + "://" + template.Host

	escapePathParmaeterReplacer := ""
	escapeParameter := false
	//replace parameters in escape path if any
	escapePathParmaeters := strings.Split(escapePath, "/")
	for _, parameter := range escapePathParmaeters {
		if regEscapeString.MatchString(parameter) {
			escapeParameter = true
			key := strings.Replace(parameter, "%7B", "", -1)
			key = strings.Replace(key, "%7D", "", -1)
			escapePathParmaeterReplacer = escapePathParmaeterReplacer + params[key] + "/"
		} else {
			escapePathParmaeterReplacer = escapePathParmaeterReplacer + parameter + "/"
		}
	}

	if escapeParameter {
		finalTemplate = finalTemplate + escapePathParmaeterReplacer + "?"
	} else {
		finalTemplate = finalTemplate + escapePath + "?"

	}
	//put parameters to template
	i := 0
	for key, value := range template.Query() {
		//check value has to be replaced or not with regexp
		if reg.MatchString(value[0]) && strings.Compare("backurl", key) != 0 {
			//parameter is a template value
			//templateParameter := regSpecial.ReplaceAllString(value[0], "")
			templateParameter := value[0]
			//template parameters can be combintation of variables
			for k, v := range params {
				templateParameter = strings.Replace(templateParameter, "{"+k+"}", v, -1)
			}
			//validate template parameter having same, send blank instead
			templateParameter = regSpecial.ReplaceAllString(templateParameter, "null")
			if i == 0 {
				finalTemplate = finalTemplate + key + "=" + templateParameter
				i = 1
			} else {
				finalTemplate = finalTemplate + "&" + key + "=" + templateParameter
			}

		} else if strings.Compare("backurl", key) == 0 {
			//get absolute string and set it as value
			templateElems := strings.Split(templateUrl, "backurl")
			valueBackUrl := strings.Replace(templateElems[1], "=", "", -1)
			if i == 0 {
				finalTemplate = finalTemplate + key + "=" + valueBackUrl
				i = 1
			} else {
				finalTemplate = finalTemplate + "&" + key + "=" + valueBackUrl
			}

		} else {
			if i == 0 {
				finalTemplate = finalTemplate + key + "=" + value[0]
				i = 1
			} else {
				finalTemplate = finalTemplate + "&" + key + "=" + value[0]
			}

		}
	}

	return finalTemplate
}

func ReplaceTemplateParameters(rawUrl string, templateUrl string, transactionID string) string {

	clickUrl, _ := url.Parse(rawUrl)
	tempUrl, _ := url.Parse(templateUrl)

	params := make(map[string]string)
	for key, value := range clickUrl.Query() {
		params[key] = value[0]
	}
	//replace transactionID if needed not present in click URL
	params["transaction_id"] = transactionID
	scheme := tempUrl.Scheme
	if scheme == "" {
		scheme = "http"
		tempUrl, _ = url.Parse(scheme + "://" + templateUrl)
	}
	escapePath := tempUrl.EscapedPath()
	finalTemplate := scheme + "://" + tempUrl.Host + escapePath + "?"

	//put parameters to template
	i := 0
	for key, value := range tempUrl.Query() {
		//check value has to be replaced or not with regexp
		if reg.MatchString(value[0]) && strings.Compare("backurl", key) != 0 {
			//parameter is a template value
			//templateParameter := regSpecial.ReplaceAllString(value[0], "")
			templateParameter := value[0]
			//template parameters can be combintation of variables
			for k, v := range params {
				templateParameter = strings.Replace(templateParameter, "{"+k+"}", v, -1)
			}
			//validate template parameter having macros , send blank instead
			templateParameter = regSpecial.ReplaceAllString(templateParameter, "null")
			if i == 0 {
				finalTemplate = finalTemplate + key + "=" + templateParameter
				i = 1
			} else {
				finalTemplate = finalTemplate + "&" + key + "=" + templateParameter
			}

		} else {
			if i == 0 {
				finalTemplate = finalTemplate + key + "=" + value[0]
				i = 1
			} else {
				finalTemplate = finalTemplate + "&" + key + "=" + value[0]
			}

		}
	}

	return finalTemplate
}
