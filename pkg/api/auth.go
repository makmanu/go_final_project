package api

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type Claims struct {
    PasswordHash string `json:"pwd_hash"`
    jwt.RegisteredClaims
}

type AuthResponse struct {
    Token string `json:"token,omitempty"`
    Error string `json:"error,omitempty"`
}

var envPassword string

func init() {
    envPassword = os.Getenv("TODO_PASSWORD")
}

func generatePasswordHash(password string) string {
    hash := sha256.Sum256([]byte(password))
    return hex.EncodeToString(hash[:])
}

func createJWTToken(password string) (string, error) {
    passwordHash := generatePasswordHash(password)
    
    expirationTime := time.Now().Add(8 * time.Hour)
    
    claims := &Claims{
        PasswordHash: passwordHash,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(expirationTime),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
        },
    }
    
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    
    secretKey := []byte(passwordHash)
    
    return token.SignedString(secretKey)
}

func validateJWTToken(tokenString, currentPassword string) (bool, error) {
    currentPasswordHash := generatePasswordHash(currentPassword)
    
    claims := &Claims{}
    
    token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
        }
        return []byte(currentPasswordHash), nil
    })
    
    if err != nil {
        return false, err
    }
    
    if !token.Valid {
        return false, nil
    }
    
    if claims.PasswordHash != currentPasswordHash {
        return false, fmt.Errorf("password changed")
    }
    
    return true, nil
}

func SignInHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }
    
    var request struct {
        Password string `json:"password"`
    }
    
    if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
        response := AuthResponse{Error: "Invalid request"}
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(response)
        return
    }
    
    envPassword := os.Getenv("TODO_PASSWORD")
    
    if envPassword == "" {
        response := AuthResponse{Error: "Authentication not configured"}
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusInternalServerError)
        json.NewEncoder(w).Encode(response)
        return
    }
    
    if request.Password != envPassword {
        response := AuthResponse{Error: "Неверный пароль"}
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusUnauthorized)
        json.NewEncoder(w).Encode(response)
        return
    }
    
    token, err := createJWTToken(envPassword)
    if err != nil {
        response := AuthResponse{Error: "Failed to create token"}
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusInternalServerError)
        json.NewEncoder(w).Encode(response)
        return
    }
    
    response := AuthResponse{Token: token}
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if envPassword == "" {
            next.ServeHTTP(w, r)
            return
        }
        
        var jwtToken string
        
        cookie, err := r.Cookie("token")
        if err == nil {
            jwtToken = cookie.Value
        } else {
            authHeader := r.Header.Get("Authorization")
            if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
                jwtToken = strings.TrimPrefix(authHeader, "Bearer ")
            }
        }
		log.Println("get token:", jwtToken)
        
        if jwtToken == "" {
            http.Error(w, "Authentication required", http.StatusUnauthorized)
            return
        }
        
        valid, err := validateJWTToken(jwtToken, envPassword)
        if err != nil || !valid {
            http.Error(w, "Authentication required", http.StatusUnauthorized)
            return
        }
        
        next(w, r)
    })
}