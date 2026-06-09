package security

import "golang.org/x/crypto/bcrypt"

// HashPassword 使用 bcrypt 对密码进行哈希处理。
func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// CheckPassword 比较明文密码与哈希值是否匹配。
func CheckPassword(hash string, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}
