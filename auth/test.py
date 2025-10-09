from sqlalchemy import create_engine, text

def test_connection_with_data():
    DATABASE_URL = "postgresql://postgres.ejajafrqufthpzctbvsn:Zov4ik2281337@aws-1-eu-west-1.pooler.supabase.com:5432/postgres"
    
    try:
        print("🔄 Пытаюсь подключиться к Supabase...")
        engine = create_engine(DATABASE_URL)

        with engine.connect() as conn:
            # Вариант 1: Получить все строки сразу
            result = conn.execute(text('SELECT * FROM "Users"'))
            all_users = result.fetchall()
            print(f"✅ Найдено пользователей: {len(all_users)}")
            
            # Вариант 2: Итерироваться по строкам (аналог курсора)
            result = conn.execute(text('SELECT * FROM "Users"'))
            print("📝 Данные пользователей:")
            for row in result.fetchall():
                print(row)
            
            return True
            
    except Exception as e:
        print(f"❌ Ошибка подключения: {e}")
        return False
test_connection_with_data()