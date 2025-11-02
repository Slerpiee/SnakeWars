import bcrypt
import jwt
from datetime import datetime, timedelta, timezone
import os
import dotenv

# Загружаем переменные окружения
dotenv.load_dotenv()

SECRET_KEY = os.getenv("JWT_KEY")
ALGORITHM = "HS256"

MIN_LOGIN_LENGTH = 3
MAX_LOGIN_LENGTH = 20
MIN_PASSWORD_LEN = 3
MAX_PASSWORD_LEN = 20
MIN_USERNAME_LEN = 3
MAX_USERNAME_LEN = 20

def validate_login(login: str) -> bool:
    if len(login) < MIN_LOGIN_LENGTH:
        return False
    if len(login) > MAX_LOGIN_LENGTH:
        return False
    if not login.isascii():
        return False
    return True

def validate_password(password: str) -> bool:
    if len(password) < MIN_PASSWORD_LEN:
        return False
    if len(password) > MAX_PASSWORD_LEN:
        return False
    if not password.isascii():
        return False
    return True

def validate_username(username: str) -> bool:
    if len(username) < MIN_USERNAME_LEN:
        return False
    if len(username) > MAX_USERNAME_LEN:
        return False
    if not username.isascii():
        return False
    return True

def generate_session_token(user_id: str, user_name: str):
    payload = {
        "user_id": user_id,
        "user_name": user_name,
        "exp": datetime.now(timezone.utc) + timedelta(hours=1)
    }
    
    token = jwt.encode(payload, algorithm=ALGORITHM, key=SECRET_KEY) 
    return token

def hash_password(password: str) -> str:
    salt = bcrypt.gensalt()
    hashed = bcrypt.hashpw(password.encode('utf-8'), salt)
    return hashed.decode('utf-8')  # Используем utf-8 для декодирования

def verify_password(plain_password: str, hashed_password: str) -> bool:
    try:
        return bcrypt.checkpw(
            plain_password.encode('utf-8'),
            hashed_password.encode('utf-8')  # Используем utf-8 для кодирования
        )
    except Exception as e:
        print(f"Error verifying password: {e}")
        return False