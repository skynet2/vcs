// Package spec provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.11.0 DO NOT EDIT.
package spec

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"net/url"
	"path"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	externalRef0 "github.com/trustbloc/vcs/pkg/restapi/v1/common"
)

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"H4sIAAAAAAAC/+x9e3PbOJL4V0Hx96uapEqyvZmZ3VvfP+uRPDu6S2Kv7XjrapNywWRLwpoCOABoRZPy",
	"d7/CiwRJgKT8yMzc7l9xROLV7250N78kKdsUjAKVIjn+koh0DRus/zxJUxDiit0BvQBRMCpA/ZyBSDkp",
	"JGE0OU7esQxytGQcmdeRfh+5AQfJJCk4K4BLAnpWrF+7keq17nRXa0DmDaTfQESIEjJ0u0NSPSrlmnHy",
	"C1avIwH8HrhaQu4KSI4TITmhq+RhkqQ3lNE0sN9L/QpKGZWYUPUnRvpVJBm6BVQKyNSfKQcsAWFUcMaW",
	"iC1RwYQAIdTCbInuYIc2WAInOEfbNVDE4ecShDRTphwyoJLgvG97N/C5IBzEDQmAYkElrICjDCjTsyoA",
	"5GQJkmwAEXX8lNFMqN2oR3ZObz1iZlAL9i101T+vj47w5ByWHMS6D6f2FTPLBG3XJF2jFFMf5OxWoQRR",
	"2DbWFEEIipQVAfSenV8tzt6fvJ0gskREoyDFuZpdHUUPcoiqqSrNCVD5n4jJNfAtETBBF6d/+7C4OJ0H",
	"19bbujE/hw6rnjjo+VQcmExD7+eScMiS4380maOx0KdJIonM1dgQX1YTs9t/QiqTSfJ5KvFKqEkZydLv",
	"7tPk08MkmVV0eWGg0MfR5xwKzAHVgxzsumydkaw71XwxV6A32LacXHMGWmNhue6WlTQLwnrJ+AYHdvmj",
	"/r2i/XrSW1C8EgX4JJGfb0J7veKYCpxqwbKI4D2McQ/b4/fRQrzZlH3Lw3UXX32Y1qvxDqKHZXcQ0zEJ",
	"Xh9S/Y9ROFsmx//40oHXl9ZOHz7pQ0u+C4jlNSvzTJGCAKmI4uriw6ni4TTHZIMyLDEiAlEm0Q4kwveY",
	"5Pg21/JKAd6cHJ2dO6WArtZEoC3Jc0RoRlIlzK9nl+hsMZ+p+fUALWcKiZUotoxKhQScKXymXVigV3NY",
	"AueQ+ZBa5mz72kPxLWM5YNrBsQc3B4cIpkcwdRDVlxLLUnSh621W6Fe6SBXV0CjZ91OwnaCPhC8br+x7",
	"rLMiIAbO9B9CE7Eaq3V7Q/c2jznuLENHUFsZeYo5yWaMLskqLB3Nsx5m+4uyVOBz4Oj2QVBS5YTeQXaT",
	"kSxADeccBFBpbChC0T+34pUZ+hoxjv4pGM2zV+ZYr5ERwAprRMJGPILl7U+Yc7zToHZgrWAziuIzuMcF",
	"0UA9/ZyuMV3BiW8OzlgGI5QamLHaTCvlGqUsA7TkbONkCFM/d/DAihtF4SOIp3rTI6DBDY+kpp55YkLe",
	"PUGbp4KgUpzd8484534y7SfAuVzP1pDe7XW0tR6HUjUwKuvSknOg8opsApPOzEOkrWGr1mtHw5kjSYYl",
	"TNU7QeM0IoeN9FCK7GMiSm3GfUyUljMLqAdlgTDNEC+pck6GDQe7lEdrIdD1Qd2ATENMg35BiSRYgtKU",
	"313PRrCUG2GUq6dsFkKUSrmii5jZ2HDobjKQmOQhBVYKyTbkFxBou8YS3RGqdbT1ExaGbLeYSm1Rrsi9",
	"9ueuZ5dh90sZFTfKqAgRlQGus4umjpaV0leM8mPOtge1LXEJ/J6kypSQAmGhTBA1covzXBkqRZEr40Mp",
	"qQ4Wqp0AzQpGaADIM23+uOeOHu15Nctu18AbhqdnMdUGdu3z4KUEjiz1Lcs83yGcqiNrcTDodxlf6YZY",
	"lN8Qi+Kbkufd7X+4eOv27F5EdqhSPv65MPq7BtkBusJ3IFDBIVVnSgExxRt24S3k+R1l28rfRgXmeANS",
	"2XyLJbpliv17NqnZqzOZMn6VcVlwdk8y5Y8al94KGDdTfQp1Mm1d2kgCSjWJRt4k1OpRxAqgJJu616bu",
	"tePDwz54VzsdE9EwtHe4ZnkG3CdBQ7FmSlQfPtVKuOTmnQ8Xb8M7qUjsRsKmyDVgQ46UfRjwiQwtWkdw",
	"uyY5NAkxZTTNy8wEO4hAhErgxic7qFxy7dqriQvOlmoKIqoTmEBCqVRCmUtS5M3l7c7ClL3imMqIV28Z",
	"LsXUUYjDtx6lPX6B5JqzcrU2e/fIUrsi9YseW+rQhwGEr5RpMwam5Ekz8qU1NqFInYYjIaEQmvq7JJzB",
	"Epe5dqmaslZNEYSDb+kEKe0e5yXYAFkVQ2lJfUV3SkQX+OcSXPjFMDiSSoIrVWfjPrdKmGs1W95OheJq",
	"KvVmTfRGH9gx+5bIdWQ9dUJkjWLtRZYFykq944LDPWGl8CBVx32QEjTkHgTC9mgK3k0cThCR6N2HyytE",
	"NIWC+r91PUV56zZ90ty01Xru+AEQCf3AQbxez2zkwCz5/uyqohVCUcO2MipJeaAmCFlwmOJKYd0YOtHC",
	"VJn8QXw7IRch/ZmRK6IWhpqGLRL1MeBzAUr7KZ1o2c/QdAFciT2FAi15mkTsYklobmhUM0U7zDgY8av2",
	"p5+LcRvzY7VdxlL4r7Voc39Gfh/47lDEYXYOzyQpBfCbgtCb2nB7pNXxgwktWDoVBaRkudMifw1yrZig",
	"joW4wxe+u6cVrdoPOl+8RzhnaqzjKRe6N1RLqI5++/RkwaO2Eg53OBs0YkOONP47o4cDWGNM0JiLPdKW",
	"0buZLSIa3pPnViLWYqnAQrFPDvdKBRBqTAGFjpZgZIHJNbTRZVkUjEth7Jufrq7O0V9Pr7SM1f+5gIxw",
	"SOWBXVagDd451KO/XRjMeTaCE6jaTlQQVEShKVwoLadNS7kGwtGG3SqW+Xtl0O4VVW2AxYk9zyg2zMY4",
	"h9wGJJaIAmRjgqZhxLm9fOohx/2c0ebw87n1G5pU5LPZHJZ6b4wusqCEKEpeMAH9znRo2YAL0XrN55Ye",
	"X8tz6wK4XMyH4xzB6ezgKOzPo7BXJ1Egr3l3HvTPap63cq4v3NcMVUeDomODYGqBPYNfzAQo1fL/n8My",
	"OU7+32F943por1sPW0e3Yc2+ELIP4gDcxpJ2eN0nx1nTtRIXdBVS2GucY7rSdgnOMmMDWnueLWOuh7I1",
	"w7dtmedrmCmUfcc2RCrzVOyEhI2J5mh/zYqjARenjqT3YS0UF36YJBnb4NBF61z/vse574GTpZWU70Cu",
	"WQQEHy4WDgLdIUb6Gps2BKEl4UIiyN58//0f/oyK8jYnqb7eZks0X8zRKyu1tWFiPK75Yv56CJoPUfp0",
	"RDaSRM+UKXTuLGMRkwjailHbNrZRgQkXvsKrbGvje5Ukz6w7zDiELVP06uLH2R//9N2fXxsbw94xqUHW",
	"yTLq3Vi5LuSjravmfNr3CzCJDVOEdYR9KiDlEJbTHcs9bjOPNVZb8qa5wsTbcXt/bi1PLLURNxLf9i5S",
	"R76UKDuJqK3Bu0xtwasZWi7T/sFIKwMOlAzYMHqww5s8KBAaC83tBC2fel8H7FrTs/YE6+vRj4ky6D4m",
	"/Z7SM2E9dKMxCkvPg/E9rq/jKI+mIzVwHg/EGub/RrTYv8nnbngQK82VeE3IfRqmzUPadhRryG6C0+1/",
	"gPOTi/5tj0qX0BlT1ogHVBYp23R9a/+Scw/TsgLVJIasgMU/jqTG0mcp1iHOHiOMSrFu0aIdHE+l+VXE",
	"UOwScxLZjg/rAfDsAWXI9ud9PWw0v/clwp2g/7o8e49oubnVoRQsEQfr14lm+p21tpyhoAwwL3MOC4RR",
	"wQSR5B6QzZQ7QFetEXXSnUBY6gkzIpQ6taGaWMojui2lsb/kriApzvOdueRRPvQ95Dsk1oxL9AoOVgcT",
	"dAtyC0DR9zpe8MejI7fR17F8PiNMSk5i2Xz1ITTbK2ibsD0LbLq6qWFCQmbj4RpkCk6C0FUO01LoLEHg",
	"YJMxDXxFAamGYiNg0Q29hkOLg5LGP2ojS7JF3zHCHJt3dykZf1S6gpCM73tRr14Lav1H8b+ezQNH/1FG",
	"Mntskj2u+h8DmREpDAM7G3m+D4VyS9v+YRTfva9XpC8kL1NponVqgDr99Sye5lBNFwznPN3d7QlPLOZJ",
	"YH6PivoBNBLK1zgnaprzGmOQjWSsezPWXg91gtxKUhaEdqEauGEJplih1oxB48rdDAQoQnmRdgP9rBra",
	"kAfoQRg9HdbDSvqxwI7fRZ4Vmu4h7gCJkILPidDZBq1d2AGTx3pMnri0U70EBoCT5a5mGZfiU+YBQjcv",
	"B5NYvaDIEpO85GDzpaxtF4rmQXoXiuSpUfqYQTQA54x3h52qn9EGhMAreHTc69p7B230S8O8Yg7idhZc",
	"yEdcD8D7cGZmjWBtKKrtYczf3e84tt2GwH7B7SD8Hg39UQHu+zbvvHR8+5kCxg9xqI2JufYCboyUryRM",
	"I3teDNGx4irRCFftQ00+U/YlG0cPtCdI/OzpMRK4cf3+u5HBvXKzw50xmDwBtENisgHWfgLbS0z5e6gE",
	"1aRxu/tMGfV7C9yu3VdvqRcljxGZITiMEZr+rvYWm/rRb0Buhg7/BPjtKzv3oO1HCc8Yuw6Lz+CpRkJG",
	"zUbokpkIBZU41UITNpjkyXGyhjxnf5G8FPI2Z+lBBvfJJKF4AzraXAr5Q85SJAFvFBh0ik6ylrIQx4eH",
	"zWGKMlrBajf8enaJhMml8S2tKttGeSI+xFEplIz++7czdD2bnpwv/JypswLoYv7dtQ6nS5YyPy/h0B3d",
	"z2I142yiczJJcpKCJQx70pMCp2uYvjk46hxyu90eYP34gPHVoR0rDt8uZqfvL0/VmAP52aDRxxrRhXKe",
	"SehS1l9dzy5fGytYGEAdHaiFtWkHFBckOU6+PTjSeymwXGv6OvQLBo6/JCuQoXiNLDkVLl4ZKctQlIxd",
	"mkzyV5A/eVPXF2N62TdHR45ywNxveGlEh0rE1uXrQ2wQKpHQ9NkSb/+tWUCUmw3mu6q0As3s/sIVFA+T",
	"5NCSgId5cWgTfGu/U+986gIIBQvFK1xBTTBNsR3uqm51urAdUYBko0w/sGz3bIAeXPbh4eHhBRE9XI80",
	"Bu2PQ4JHIJUTH6ONwlxXTXXW6TTDEmsq+WXqXROGCcRedAmkbwrDN91+7oOXONW4COySjJ05crH7EtQy",
	"6k75hSlm3MXhGKoZm4fwKDpp+P5hyvhgky6rSiVP31WFN5JV0eVmnYYtxbAJtc1s1BipzJoFzs9PIN1q",
	"+BemhoC7uBfqfYCMR3Ip1i0lMSgGOsi2eVz+9b5OgddhdyS923ttkjQlm+eEthAduep9KYEwcLMcx/4Q",
	"gqLX8vsgSkjG91Pn+nJJPFWZD93AvQQq+td8YTYcuJMbw5KPgfw+tGBvOGDavFkYoAd3NSCi1yKldw/U",
	"pIIR1wovQQiDy74wLQzfQ40hh/GAHyACW4QoDr/Yvxbzh0Mv+Gne0xTgpbD+I1Jg6BJsTQMaop4ov6t2",
	"hKtFEj8YJHkJEw9+7fjGp0mE/hbtFPqQzGdCtjJ4X0rehxLZn4GcWpG0EeShN4LSsep7kAjqhgC/RSow",
	"d/HCtwBjProiBo8Oqo4uL0EN/SkCvwpd9ELqGSjk8Iv5dzF/6AutcAL3INrV1T1xlRDKfkVKnIQbY+hZ",
	"AouI+ule1P6ViWMEYvYmkYaNUTUNYCRLf7PCxCuBJFUJJPHrMxfB0KsfJiVU9+6wGZPNOzMR68YQ6ilY",
	"vaq2o8ukfXPPL0Ts8o2rXKuZxxWTvpT6C1fNvrBBFSuOHKUnh8puYzSvKPiwMruiss4Etd4cHIVDJ66R",
	"pEWG6f2i+yxU/RPalfd+rnsT3WckSyuLckg+DhZlaEb7uQS+qzmtXVfxBIF5VbcvMF0ul8R4iqF1/aqd",
	"J6x5gqrrO5QBJ/eQVXXCxm6urjxciwddfWyTdYMZuhNbOG1HZgivFItL09kieiCWwU19l/jEU5nkIbPn",
	"La77Upgz2sppt9i4Ld2YOZO9cRrM9ua2uNsoFuWYTPEKaNVYwuD3G1G92Oit45pf5DsEQuLbnOiE+arq",
	"P7ikbZzR6JKxIkLaHjEFZ5q/GDdtJzb4zr0eTcQOc4TZsM2/3hNYpp1rs03twIKmNGo/AqGujYmpLPRr",
	"/C1sJEMbTEwfINPJw6Xc+0UCuvEQzvNbnN4Z1RYEve0wIkwLErOmbTBssWsh7RGCmrJJDWaBuqHI5U9n",
	"H97OK9VoUyjulejQpbxMiKkgst7tkvEV8F0UkFW64+Pp2xWTKM1+DztD3u43fMtK2bKkzBu2zrLq9mX6",
	"Bh+gd67lT2QRzzIwxK8bt+p+HjfN4HeFsQZ+CEUpNjf0ge5CIgapcP3MXpAzF8bfCHvjjGaMUkilq0j9",
	"cPHWoNv1QSN5rpuduIoSdg98VzGtFm0S+IZQ8AD6jQJRgW9JTiQBocnVCRFxgC5OZ2fv3p2+n5/OFSTm",
	"O4o3JPVV60U/65lVbqwZ8EgW1FGZtQ5m15Tw7uR/9HEV99UVIY7VbNsWSTbkF6gY5xuhG6RwAjSFZzid",
	"zuRem2q0vRwfr6WS1eQ720IcuBYoFm2uxxd8lq5YqGVFAz9AJ9EWRkod19VCBRa2nRCmwdZslRhwCr62",
	"5WvI21KeTic2v7uTbn6ihtRtjswWGzKre5Kres1NKSSS+E77CExJelZS20eqmtR2B16VWBmAYFuKc7Ii",
	"VD225yDCTjpBqes4jCnCUiqhHMGtl939eP/z26M3PSb75+l2u50uGd9MS54DVeZE1rThwzU+rdiAa1se",
	"UC/ajlkBVTbvQEf92Ght75oiKVNhlu9sP0GizT3b2U2pQyLJyvldnIg7JTVzwHeRbu7hJH93HNfx7aN5",
	"8WPikZqy2FyfG2tpWq0c6TalzgafcSotHdr2X74tazTocFakK7UYChX8yEqatVwn7TENXQjXJWSV01Rg",
	"Hr9PmJmTC6CZywgIl+IZSyLfdfrCOStESf8VSNEucay7gCnu83UqFt36PVes54llXjdFixdPdz2zYBHe",
	"fvcce/PYyOak/wJ2TbQNaKQhSdA17U7SdOOOfxsO58A2nWt3/AyO5GN7L/7bUPj1DYVAGw/Plz3+F3Pu",
	"v2JTk73jAGOtjX87+uH61HWwn8hvzCfrbL3pbh7/7l3qoT4GPY2ammo2ZLB27xf+8KxJXLH2CQFbeWZ7",
	"vT1Mku+Ovg+UTBgl+55JdJLnbGtf/cO34aaqisJPqSRyh64YQ28xX4Ee8ObPoW6VDL3DdOfgLkI2e6Th",
	"yBjTfZ/UoFMn4vVnDTDNcsW/rdawVT/WsM3cTpF5OVO5r3m21sq3gODnEueo5PSYgFwea5EljnXK1bEe",
	"P1Xjj0P17AH+HlWHP3KueBl+9WSS0DLXnxxyrn8/Q3oAmURL9IfZ8Pmu+UJfDYtc8SnOC2RR/oAzVN9J",
	"dvkiQpqGeBvZdAHecDovfhWo1nBvNZmi8km9u91upy6llJnSgiUgVvJmR+pYZrNiI9dbd+hO0GtEUifZ",
	"e7mEsWujp91fuchUX5D+KVGroFy1AAlQgQesHmxXXwrsSUsTzSb4pp28ekH7Jbj5ScGqv6WrstGdrrsR",
	"LBO1MMbcGgvrYwc+i9HTu6xLIVf2e30vJFx7nO3ORxGc522cHL+Lvt/pO/YdzQEZN4mEHSPhwm5kbfz3",
	"O5XNUhUPhiIDfFdItuK4WFtPlWOasQ0yc3Q+MlB3p4v3p7J2tSGoPgeiT9tFPJvuxxcifk6/G9nB98fG",
	"gE5k1PoY2UBMBJsvMBD74QFR7b/6hOawCdrQeFauNdH4f07lNRrmRyWehfLUHPrwS1mS7GGwctIRpRnV",
	"lTt21TP9+Ifdh9Kmd+xdC9HuMmUWVI5vaeYMfICtVw+qYUr8NicMJ5+V5Z5ZKWr1qrS2mR7XLpnx+qBb",
	"LdOEoC7jjLV1eyFRTrLYV3IXc0tOWiqaHoK0+xmUFEhhIgxV+GADEutAex1avj43k+0T/riUVfgrLB5b",
	"nfqDjXuK2PHcjuptMwqIcbRhHJBXp+yXmItIsf5IIdI6X6mkgdrl96HHP5o2HO2EXxsDqgLNjW8XNBIb",
	"TRyhUbuvzV/3VaDr2aXHTH5hfJSiv8jPOk84x2TjCYy2IDBppwtvpK45fGoFdRM+ujWm+SqK17zEz3It",
	"HZ8/Iq96CMwrkGZxz2SzYUwjbIvGx6LCgB7KeZ7rGGJdGhcWWfoDB48XWYOZ4KaRxnB679y0itBzvFBm",
	"bzf3v91F56Vy/4Ndn166sijWIWhUQVG7Z9QIXn/2jO+vThJV7jDJUk/+fJX86POvQROxr5c8i1B7bt0R",
	"pCd/0t+FcPENgBeVLp0WSV9FvgRb6OwhYYomeCI04SjgalfAQ5gwtpDnU/0ByMOMZNO0+kp1r1tSv9p1",
	"SepvXb8gFOtFxtXkuHyX6oT7OzSuhObK+LZxZrh6cl5/Va2TPSvXaaDoaIs5X91Y6PjwMGcpztdMyOP/",
	"OPrTUaJ41EKovTsT05yaQEpmejG3brXqrdqbt+4ZHamOnKei7EDss9tdqB7nd+V5+PTwvwEAAP//Uygh",
	"DOmGAAA=",
}

// GetSwagger returns the content of the embedded swagger specification file
// or error if failed to decode
func decodeSpec() ([]byte, error) {
	zipped, err := base64.StdEncoding.DecodeString(strings.Join(swaggerSpec, ""))
	if err != nil {
		return nil, fmt.Errorf("error base64 decoding spec: %s", err)
	}
	zr, err := gzip.NewReader(bytes.NewReader(zipped))
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %s", err)
	}
	var buf bytes.Buffer
	_, err = buf.ReadFrom(zr)
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %s", err)
	}

	return buf.Bytes(), nil
}

var rawSpec = decodeSpecCached()

// a naive cached of a decoded swagger spec
func decodeSpecCached() func() ([]byte, error) {
	data, err := decodeSpec()
	return func() ([]byte, error) {
		return data, err
	}
}

// Constructs a synthetic filesystem for resolving external references when loading openapi specifications.
func PathToRawSpec(pathToFile string) map[string]func() ([]byte, error) {
	var res = make(map[string]func() ([]byte, error))
	if len(pathToFile) > 0 {
		res[pathToFile] = rawSpec
	}

	pathPrefix := path.Dir(pathToFile)

	for rawPath, rawFunc := range externalRef0.PathToRawSpec(path.Join(pathPrefix, "./common.yaml")) {
		if _, ok := res[rawPath]; ok {
			// it is not possible to compare functions in golang, so always overwrite the old value
		}
		res[rawPath] = rawFunc
	}
	return res
}

// GetSwagger returns the Swagger specification corresponding to the generated code
// in this file. The external references of Swagger specification are resolved.
// The logic of resolving external references is tightly connected to "import-mapping" feature.
// Externally referenced files must be embedded in the corresponding golang packages.
// Urls can be supported but this task was out of the scope.
func GetSwagger() (swagger *openapi3.T, err error) {
	var resolvePath = PathToRawSpec("")

	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true
	loader.ReadFromURIFunc = func(loader *openapi3.Loader, url *url.URL) ([]byte, error) {
		var pathToFile = url.String()
		pathToFile = path.Clean(pathToFile)
		getSpec, ok := resolvePath[pathToFile]
		if !ok {
			err1 := fmt.Errorf("path not found: %s", pathToFile)
			return nil, err1
		}
		return getSpec()
	}
	var specData []byte
	specData, err = rawSpec()
	if err != nil {
		return
	}
	swagger, err = loader.LoadFromData(specData)
	if err != nil {
		return
	}
	return
}
