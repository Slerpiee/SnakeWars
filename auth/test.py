# from sqlalchemy import create_engine, text
# import time

# def test_connection():
#     DATABASE_URL = "postgresql://postgres.ejajafrqufthpzctbvsn:Zov4ik2281337@aws-1-eu-west-1.pooler.supabase.com:5432/postgres"
    
#     try:
#         start_time = time.time()

#         engine = create_engine(DATABASE_URL)

#         with engine.connect() as conn:
#             result = conn.execute(text("SELECT version();"))
#             db_version = result.fetchone()[0]
            
#             end_time = time.time()
#             connection_time = round((end_time - start_time) * 1000, 2)
            
#             print(f"✅ Подключение успешно!")
#             print(f"📊 Версия PostgreSQL: {db_version}")
#             print(f"⏱️  Время подключения: {connection_time} ms")
#             print(f"🔗 Connection string: {DATABASE_URL.split('@')[1]}")
            
#             return True
            
#     except Exception as e:
#         print(f"❌ Ошибка подключения: {e}")
#         return False

# if __name__ == "__main__":
#     test_connection()
from sqlalchemy import create_engine
# from sqlalchemy.pool import NullPool
from dotenv import load_dotenv
import os
from sqlalchemy.pool import NullPool

# Load environment variables from .env
load_dotenv()

# Fetch variables
USER = os.getenv("user")
PASSWORD = os.getenv("password")
HOST = os.getenv("host")
PORT = os.getenv("port")
DBNAME = os.getenv("dbname")

# Construct the SQLAlchemy connection string
DATABASE_URL = f"postgresql+psycopg2://{USER}:{PASSWORD}@{HOST}:{PORT}/{DBNAME}?sslmode=require"

# Create the SQLAlchemy engine
engine = create_engine(DATABASE_URL, poolclass=NullPool)
# If using Transaction Pooler or Session Pooler, we want to ensure we disable SQLAlchemy client side pooling -
# https://docs.sqlalchemy.org/en/20/core/pooling.html#switching-pool-implementations
# engine = create_engine(DATABASE_URL, poolclass=NullPool)

# Test the connection
try:
    with engine.connect() as connection:
        print("Connection successful!")
except Exception as e:
    print(f"Failed to connect: {e}")