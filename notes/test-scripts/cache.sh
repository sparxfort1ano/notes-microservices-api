#!/bin/bash

AUTH_URL="http://localhost:8081/auth"
NOTES_URL="http://localhost:8082/notes"

echo "=== Получение токена ==="
TEST_USER="cachetester_$(date +%s)"

curl -s -X POST "$AUTH_URL/register" -H "Content-Type: application/json" -d "{\"username\": \"$TEST_USER\", \"password\": \"testpass123\"}" > /dev/null
LOGIN_RESPONSE=$(curl -s -X POST "$AUTH_URL/login" -H "Content-Type: application/json" -d "{\"username\": \"$TEST_USER\", \"password\": \"testpass123\"}")
ACCESS_TOKEN=$(echo "$LOGIN_RESPONSE" | grep -o '"access_token":"[^"]*"' | cut -d'"' -f4)

if [ -z "$ACCESS_TOKEN" ]; then
    echo "Ошибка авторизации."
    exit 1
fi

echo -e "Токен получен. Смотри в логи сервера!\n"
sleep 2

# 1. Создание
echo "1. POST: Создаем заметку..."
CREATE_RESPONSE=$(curl -X POST "$NOTES_URL" -H "Authorization: Bearer $ACCESS_TOKEN" -H "Content-Type: application/json" -d '{"name": "Кэш-тест", "content": "Первоначальный текст"}' -s)
NOTE_ID=$(echo "$CREATE_RESPONSE" | grep -o '"id":"[^"]*"' | cut -d'"' -f4)
sleep 2

# 2. Прогрев кэша списка (Cache Miss)
echo -e "\n2. GET ALL (1-й раз): Идем в Mongo, сохраняем в Redis. Ожидай бОльшую латенси."
curl -X GET "$NOTES_URL" -H "Authorization: Bearer $ACCESS_TOKEN" -s > /dev/null
sleep 2

# 3. Чтение из кэша списка (Cache Hit)
echo -e "\n3. GET ALL (2-й раз): Читаем из Redis. Латенси должна резко упасть!"
curl -X GET "$NOTES_URL" -H "Authorization: Bearer $ACCESS_TOKEN" -s > /dev/null
sleep 2

# 4. Прогрев кэша одной заметки (Cache Miss)
echo -e "\n4. GET /id (1-й раз): Идем в Mongo, сохраняем в Redis."
curl -X GET "$NOTES_URL/$NOTE_ID" -H "Authorization: Bearer $ACCESS_TOKEN" -s > /dev/null
sleep 2

# 5. Чтение из кэша одной заметки (Cache Hit)
echo -e "\n5. GET /id (2-й раз): Читаем из Redis. Сравни латенси с шагом 4!"
curl -X GET "$NOTES_URL/$NOTE_ID" -H "Authorization: Bearer $ACCESS_TOKEN" -s > /dev/null
sleep 2

# 6. Инвалидация
echo -e "\n6. PATCH: Обновляем заметку. Это снесет ключи в Redis (InvalidateNotes)."
curl -X PATCH "$NOTES_URL/$NOTE_ID" -H "Authorization: Bearer $ACCESS_TOKEN" -H "Content-Type: application/json" -d '{"content": "Обновленный текст"}' -s > /dev/null
sleep 2

# 7. Проверка инвалидации (Cache Miss)
echo -e "\n7. GET /id (после PATCH): Снова идем в Mongo, т.к. кэш был сброшен. Данные должны быть свежими."
curl -X GET "$NOTES_URL/$NOTE_ID" -H "Authorization: Bearer $ACCESS_TOKEN" -s
echo ""
sleep 2

# 8. Удаление
echo -e "\n8. DELETE: Удаляем заметку (чистит базу и кэш)."
curl -X DELETE "$NOTES_URL/$NOTE_ID" -H "Authorization: Bearer $ACCESS_TOKEN" -s > /dev/null

echo -e "\nТестирование кэша завершено!"