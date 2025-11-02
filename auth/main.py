from pydantic import BaseModel
from fastapi import FastAPI, Response, HTTPException, Depends
from verify import *
from database import *
import sys

#ЗАПУСК: uvicorn main:app --host 0.0.0.0 --port 8000 --reload

app = FastAPI()

@app.on_event("startup")
async def startup_event():
    print("starting...")
    try:
        await check_db_connection()
    except Exception as e:
        print("Failed to connect to database: ", e)
        sys.exit(-1)


class LoginRequest(BaseModel):
    login: str
    password: str

class RegisterRequest(BaseModel):
    login: str
    password: str
    username: str



@app.get("/")
async def root():
    return {"api is working"}


@app.post("/register")
async def login(register_data: RegisterRequest, db: AsyncSession = Depends(get_db)):

    print(register_data, "<----")

    if not(validate_login(register_data.login)):
        return HTTPException(status=400)
    
    try:
        db_res = await db.execute(text("SELECT login FROM Users WHERE login=:user_login"), {"user_login":register_data.login}) 
    except Exception:
        return HTTPException(status_code=503)
    
    if db_res.scalar_one_or_none() is not None:
            return HTTPException(status_code=400, detail="User already registered")
    if not(validate_password(register_data.password)):
        return HTTPException(status_code=400)
    
    if not(validate_username(register_data.username)):
        return HTTPException(status_code=400)
    register_data.password = hash_password(register_data.password)
    try:
        await db.execute(text("INSERT INTO Users (login, hashed_password, username) VALUES (:login, :password, :username)"), {"login": register_data.login, "password": register_data.password, "username": register_data.username}) #Добавить отлавливание ошибки и возврат 5-- кода ошибки
    except Exception:
        return HTTPException(status_code=503)
    return {"success"} #после удачной реги, пользователь должен будет залогиниться 

@app.post("/login")
async def login(login_data: LoginRequest, db: AsyncSession = Depends(get_db)):
    if not(validate_login(login_data.login)):
        return HTTPException(status_code=403)
     
    if not(validate_password(login_data.password)):
        return HTTPException(status_code=403)
    try:
        db_res = await db.execute(text("SELECT id, username, hashed_password FROM Users WHERE login = :login"), {"login": login_data.login}) 
    except Exception:
        return HTTPException(status_code = 503)
    user = db_res.first()
    if user is None:
        return HTTPException(status_code=403, detail="User is not registered")
    user_data = user._asdict()
    print(user_data)
    if not verify_password(login_data.password, user_data["hashed_password"]):
        return HTTPException(status_code=403, detail="Incorrect password or login")
    jwt_token = generate_session_token(user_id=user_data["id"], user_name=user_data["username"])

    return {
        "access_token": jwt_token,
        "token_type": "bearer" 
    }
    
    
    
    

    

    

