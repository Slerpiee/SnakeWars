from sqlalchemy import create_engine, text
import time

def test_connection():
    DATABASE_URL = "postgresql://postgres.ejajafrqufthpzctbvsn:Zov4ik2281337@aws-1-eu-west-1.pooler.supabase.com:5432/postgres"
    
    try:
        print("🔄 Пытаюсь подключиться к Supabase...")
        start_time = time.time()

        engine = create_engine(DATABASE_URL)

        with engine.connect() as conn:
            result = conn.execute(text("SELECT version();"))
            db_version = result.fetchone()[0]
            
            end_time = time.time()
            connection_time = round((end_time - start_time) * 1000, 2)
            
            print(f"✅ Подключение успешно!")
            print(f"📊 Версия PostgreSQL: {db_version}")
            print(f"⏱️  Время подключения: {connection_time} ms")
            print(f"🔗 Connection string: {DATABASE_URL.split('@')[1]}")
            
            return True
            
    except Exception as e:
        print(f"❌ Ошибка подключения: {e}")
        return False

# Запуск теста
if __name__ == "__main__":
    test_connection()