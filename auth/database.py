from sqlalchemy.ext.asyncio import *
from sqlalchemy.orm import * 
from sqlalchemy.pool import *
from sqlalchemy import text

Base = declarative_base()

DATABASE_URL = "postgresql://postgres.ejajafrqufthpzctbvsn:Zov4ik2281337@aws-1-eu-west-1.pooler.supabase.com:5432/postgres"
    

engine = create_async_engine(
    DATABASE_URL,
    echo=True,  
    poolclass=NullPool,  
    future=True,  
    pool_pre_ping=True, 
)


AsyncSessionLocal = sessionmaker(
    engine,
    class_=AsyncSession,
    expire_on_commit=False, 
)   

async def get_db():
    async with AsyncSessionLocal() as session:
        try:
            yield session # выдаем сессию хэндлеру и потом автоматически коммитем все изменения
            await session.commit()  
        except Exception:
            await session.rollback() 
            raise
        finally:
            await session.close()


# 1. Запрос → FastAPI
# 2. FastAPI вызывает get_db() → создает генератор
# 3. generator.__anext__() → входит в async with (создается сессия)
# 4. generator.__anext__() → доходит до yield session
# 5. session передается в ваш эндпоинт
# 6. Ваш эндпоинт выполняется (работа с БД)
# 7. После завершения эндпоинта → generator.__anext__() продолжается
# 8. Выполняется await session.commit()
# 9. Если была ошибка → await session.rollback()
# 10. finally → await session.close()
# 11. Выход из async with
# 12. Ответ отправляется клиенту


async def check_db_connection():
    try:
        async with engine.begin() as conn:
            await conn.execute(text("SELECT 1"))
        return True
    except Exception as e:
        print(f"Database connection error: {e}")
        return False

async def close_db():
    await engine.dispose()