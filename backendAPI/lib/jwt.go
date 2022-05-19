package lib

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"time"

	//"fmt"

	//"fmt"
	"net/http"

	//"io/ioutil"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

var (
	TokenExpired     = errors.New("token過期")
	TokenNotValidYet = errors.New("token不合法")
	TokenMalformed   = errors.New("token不可用")
	TokenInvalid     = errors.New("無效的token")
)

var (
	Manager = 2
	Editor  = 1
	Viewer  = 0
)

// 定義一個JWTAuth的中介軟體
// func JWTAuth() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		// 通過http header中的token解析來認證
// 		//token := c.Request.Header.Get("token")
// 		reqToken := c.Request.Header.Get("Authorization")
// 		splitToken := strings.Split(reqToken, "Bearer ")
// 		token := splitToken[1]
// 		if token == "" {
// 			//c.String(http.StatusUnauthorized, "請求未攜帶token，無許可權訪問")
// 			c.JSON(http.StatusForbidden, gin.H{
// 				"status": -1,
// 				"msg":    "請求未攜帶token，無許可權訪問",
// 			})
// 			c.Abort()
// 			return
// 		}

// 		fmt.Println("get token: ", token)

// 		// A, _ := jwt.ParseECPublicKeyFromPEM([]byte(`-----BEGIN PUBLIC KEY-----
// 		// MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEEVs/o5+uQbTjL3chynL4wXgUg2R9
// 		// q9UU8I5mEovUf86QZ7kOBIjJwqnzD1omageEHWwHdBO6B+dFabmdT9POxg==
// 		// -----END PUBLIC KEY-----`))

// 		// 初始化一個JWT物件例項，並根據結構體方法來解析token
// 		// 解析token中包含的相關資訊(有效載荷)
// 		claims, err := ParseToken(token)

// 		if err != nil {
// 			// token過期
// 			if err == TokenExpired {
// 				c.JSON(http.StatusUnauthorized, gin.H{
// 					"status": -1,
// 					"msg":    "token授權已過期，請重新申請授權",
// 				})
// 				c.Abort()
// 				return
// 			}
// 			// 其他錯誤
// 			c.JSON(http.StatusBadRequest, gin.H{
// 				"status": -1,
// 				"msg":    err.Error(),
// 			})
// 			c.Abort()
// 			return
// 		}
// 		// 將解析後的有效載荷claims重新寫入gin.Context引用物件中
// 		c.Set("claims", claims)
// 	}
// }
//var privateKey *ecdsa.PrivateKey

// var publicKey, _ = jwt.ParseECPublicKeyFromPEM([]byte(`-----BEGIN PUBLIC KEY-----
//   MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEEVs/o5+uQbTjL3chynL4wXgUg2R9
//   q9UU8I5mEovUf86QZ7kOBIjJwqnzD1omageEHWwHdBO6B+dFabmdT9POxg==
//   -----END PUBLIC KEY-----`))

type Claims struct {
	Act string `json:"act"`
	Dom int    `json:"dom"`
	jwt.StandardClaims
}

type MyClaims struct {
	Act string
	Dom int
}

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		//驗證
		auth(c)

		//處理請求
		c.Next()

		endTime := time.Now()

		// 執行時間
		latencyTime := endTime.Sub(startTime)

		//日誌格式
		Log.Infof(" %3d | %13v | %15s | %s | %s ",
			c.Writer.Status(),
			latencyTime,
			c.ClientIP(),
			c.Request.Method,
			c.Request.RequestURI,
		)
	}
}

func auth(c *gin.Context) {
	var err error
	// privateKey, _, err := getEcdsaKey(2) //生成椭圆曲线的私钥
	// if err != nil {
	// 	fmt.Println("getEcdsaKey is error!", err)
	// 	c.Abort()

	// 	return
	// }
	// token, err := ReleaseToken(privateKey)
	// if err != nil {
	// 	fmt.Println("生成token错误：", err)
	// 	return
	// }
	// fmt.Println("生成的token为：", token)

	reqToken := c.Request.Header.Get("Authorization")
	splitToken := strings.Split(reqToken, "Bearer ")
	tokenString := splitToken[1]
	var publicKey *ecdsa.PublicKey

	// keyBytes, err := ioutil.ReadFile("pub.pem")
	// if err != nil {
	// 	fmt.Println("unable to read public key")
	// }

	// publicKey, err = jwt.ParseECPublicKeyFromPEM(keyBytes)

	// if err != nil {
	// 	fmt.Println(err)
	// }

	//	token, err := ReleaseToken(privateKey)
	// if err != nil {
	// 	fmt.Println("生成token错误：", err)
	// 	return
	// }

	// type Result struct {
	// 	ID   int
	// 	Name string
	// 	Age  int
	// }

	//fmt.Printf("%s\n", tokenString)
	account, err := getAccount(tokenString)
	//fmt.Printf("%s\n", tokenString)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": -1,
			"msg":    "token格式異常",
		})
		c.Abort()
		return
	}
	var jwtKey string
	var keyExpire int
	var admin int
	var accountId int
	var customerId int
	var token string
	//var result Result
	err = DB.Raw("SELECT jwtKey,keyExpire,admin,accountID,customerID,token FROM account WHERE account = ?", account).Row().Scan(&jwtKey, &keyExpire, &admin, &accountId, &customerId, &token)
	// fmt.Printf(">a>%s  %s [%s]\n", account, err, jwtKey)

	if err != nil {

		c.JSON(http.StatusUnauthorized, gin.H{
			"status": -3,
			"msg":    "token無效,請重新登入",
		})
		c.Abort()
		return
	}
	publicKey, err = jwt.ParseECPublicKeyFromPEM([]byte(jwtKey))
	// 	publicKey, _ = jwt.ParseECPublicKeyFromPEM([]byte(`-----BEGIN PUBLIC KEY-----
	// MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEfyfmfy+LDfK9JcGJIAokOmKlXLn0
	// I8ZGR+R2AaI3d+T2daGTezplwFCESpz5XnINYIkCBXD2knOApuJpxOm6ZA==
	// -----END PUBLIC KEY-----`))
	if err != nil {
		//fmt.Printf(">b>  %s\n", err)

		c.JSON(http.StatusUnauthorized, gin.H{
			"status": -2,
			"msg":    "token無效,請重新登入",
		})
		c.Abort()
		return
	}

	//fmt.Println("生成的token为：", token)
	//parseToken, claims, err := ParseToken(token, publicKey)
	_, err = Decode(tokenString, publicKey)
	//fmt.Println(privateKey)
	//fmt.Println(parseToken, claims, err)
	//act := claims.act
	//fmt.Println(claims.act)

	if err != nil {
		//fmt.Printf(">c>  %s [%s]\n", err, jwtKey)

		c.JSON(http.StatusUnauthorized, gin.H{
			"status": -4,
			"msg":    "token無效,請重新登入",
		})
		c.Abort()
		return
	}

	c.Set("admin", admin)
	c.Set("id", accountId)
	c.Set("customer", customerId)
	c.Set("token", token)
	////// 待處理 core dump
	//c.Keys["claims"] = claims

}

func getAccount(tokenString string) (string, error) {
	var claims Claims
	parts := strings.Split(tokenString, ".")
	if len(parts) != 3 {
		return "", errors.New("token contains an invalid number of segments")
	}
	claimBytes, err := jwt.DecodeSegment(parts[1])
	if err != nil {
		return "", err
	}

	//err = json.NewDecoder(bytes.NewBuffer(claimBytes)).Decode(&claims)
	err = json.Unmarshal(claimBytes, &claims)
	//fmt.Printf(">z>_%s_%s_%s_%s\n", parts[1], claims.Act, string(claimBytes), err)
	//fmt.Println(claims)

	return claims.Act, err
}

//生成token
// func ReleaseToken(key *ecdsa.PrivateKey) (string, error) {
// 	expirationTime := time.Now().Add(7 * 24 * time.Hour) //截止时间：从当前时刻算起，7天
// 	claims := &Claims{
// 		UserId: 001, //分发给某一用户的token，模拟数据为 001
// 		StandardClaims: jwt.StandardClaims{
// 			ExpiresAt: expirationTime.Unix(), //过期时间
// 			IssuedAt:  time.Now().Unix(),     //发布时间
// 			Issuer:    "jiangzhou",           //发布者
// 			Subject:   "user token",          //主题
// 		},
// 	}
// 	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims) //生成token
// 	tokenString, err := token.SignedString(key)                //签名
// 	if err != nil {
// 		fmt.Println("生成token错误：", err)
// 		return "", err
// 	}
// 	//fmt.Println("token:",tokenString)
// 	return tokenString, err
// }

//解析token
// func ParseToken(tokenString string, key *ecdsa.PublicKey) (*jwt.Token, *Claims, error) {
// 	claims := &Claims{}
// 	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (i interface{}, err error) {
// 		//token.Method = jwt.SigningMethodES256
// 		return key, nil
// 	})
// 	return token, claims, err
// }

// func GenKey() (string, string, error) {

// 	privateKey, _, err := getEcdsaKey(2) //生成椭圆曲线的私钥
// 	x509Encoded, _ := x509.MarshalECPrivateKey(privateKey)
// 	privateBs := (string)(pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: x509Encoded}))
// 	//fmt.Println(privateBs)
// 	x509EncodedPub, _ := x509.MarshalPKIXPublicKey(privateKey.Public())
// 	publicBs := (string)(pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: x509EncodedPub}))
// 	//fmt.Println(publicBs)
// 	return privateBs, publicBs, err
// }

// func getEcdsaKey(length int) (*ecdsa.PrivateKey, *ecdsa.PublicKey, error) {
// 	var err error
// 	var prk *ecdsa.PrivateKey
// 	var puk *ecdsa.PublicKey
// 	var curve elliptic.Curve
// 	switch length {
// 	case 1:
// 		curve = elliptic.P224()
// 	case 2:
// 		curve = elliptic.P256()
// 	case 3:
// 		curve = elliptic.P384()
// 	case 4:
// 		curve = elliptic.P521()
// 	default:
// 		err = errors.New("输入的签名等级错误！")
// 	}
// 	prk, err = ecdsa.GenerateKey(curve, rand.Reader) //通过 "crypto/rand" 模块产生的随机数生成私钥
// 	puk = &prk.PublicKey

// 	//if err != nil {
// 	return prk, puk, err
// 	//}
// }

func Decode(token string, key *ecdsa.PublicKey) (*Claims, error) {
	tokenType, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return key, nil
	})

	//fmt.Printf(">g>  %s\n", to)

	if claims, ok := tokenType.Claims.(*Claims); ok && tokenType.Valid {
		return claims, nil
	} else {
		//fmt.Printf("%s\n", tokenType.Valid)
		return nil, err
	}
}

func KeyPairs() (string, string, error) {
	//elliptic.P256(),elliptic.P384(),elliptic.P521()

	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return "", "", err
	}

	//x509Encoded, _ := x509.MarshalECPrivateKey(privateKey)
	x509Encoded, _ := x509.MarshalPKCS8PrivateKey(privateKey)
	privateBs := pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: x509Encoded})

	//privateFile, err := os.Create("private.pem")
	//_, err = privateFile.Write(privateBs)

	x509EncodedPub, _ := x509.MarshalPKIXPublicKey(privateKey.Public())
	publicBs := pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: x509EncodedPub})

	//publicKeyFile, err := os.Create("public.pem")
	//_, err = publicKeyFile.Write(publicBs)

	return string(privateBs), string(publicBs[:len(string(publicBs))-1]), nil
}
