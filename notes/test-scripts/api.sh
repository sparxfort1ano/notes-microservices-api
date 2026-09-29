AUTH_URL="http://localhost:8081/auth"
NOTES_URL="http://localhost:8082/notes"

echo "=== Подготовка: Получение токена из сервиса Auth ==="
TEST_USER="notestester_$RANDOM"

curl -s -X POST "$AUTH_URL/register" \
     -H "Content-Type: application/json" \
     -d "{\"username\": \"$TEST_USER\", \"password\": \"testpass123\"}" > /dev/null

LOGIN_RESPONSE=$(curl -s -X POST "$AUTH_URL/login" \
     -H "Content-Type: application/json" \
     -d "{\"username\": \"$TEST_USER\", \"password\": \"testpass123\"}")

ACCESS_TOKEN=$(echo "$LOGIN_RESPONSE" | grep -o '"access_token":"[^"]*"' | cut -d'"' -f4)

if [ -z "$ACCESS_TOKEN" ]; then
    echo "Ошибка: Не удалось получить токен. Проверь, запущен ли сервис auth."
    exit 1
fi

echo "Успешно получен токен для $TEST_USER"
echo "----------------------------------------"
sleep 1

echo "=== Тестирование Notes API ==="
echo "----------------------------------------"

# 1. Создание новой заметки
echo ""
echo "Запрос 1: POST $NOTES_URL"

CREATE_RESPONSE=$(curl -X POST "$NOTES_URL" \
     -H "Authorization: Bearer $ACCESS_TOKEN" \
     -H "Content-Type: application/json" \
     -d '{
           "name": "Bash test note",
           "content": "Текст заметки из bash-скрипта"
         }' \
     -w "\nStatus: %{http_code}\n" \
     -s)

echo "$CREATE_RESPONSE"

NOTE_ID=$(echo "$CREATE_RESPONSE" | grep -o '"id":"[^"]*"' | cut -d'"' -f4)

echo "Извлеченный ID заметки: $NOTE_ID"
echo "----------------------------------------"
sleep 1

# 2. Получение списка всех заметок
echo ""
echo "Запрос 2: GET $NOTES_URL"

curl -X GET "$NOTES_URL" \
     -H "Authorization: Bearer $ACCESS_TOKEN" \
     -w "\nStatus: %{http_code}\n"
echo "----------------------------------------"
sleep 1

# 3. Получение конкретной заметки по ID
echo ""
echo "Запрос 3: GET $NOTES_URL/$NOTE_ID"

curl -X GET "$NOTES_URL/$NOTE_ID" \
     -H "Authorization: Bearer $ACCESS_TOKEN" \
     -w "\nStatus: %{http_code}\n"
echo "----------------------------------------"
sleep 1

# 4. Обновление заметки
echo ""
echo "Запрос 4: PATCH $NOTES_URL/$NOTE_ID"

curl -X PATCH "$NOTES_URL/$NOTE_ID" \
     -H "Authorization: Bearer $ACCESS_TOKEN" \
     -H "Content-Type: application/json" \
     -d '{
           "name": "Обновленный заголовок",
           "content": "Обновленный текст из скрипта"
         }' \
     -w "\nStatus: %{http_code}\n"
echo "----------------------------------------"
sleep 1

# 5. Проверка обновленной информации
echo ""
echo "Запрос 5: GET $NOTES_URL/$NOTE_ID"

curl -X GET "$NOTES_URL/$NOTE_ID" \
     -H "Authorization: Bearer $ACCESS_TOKEN" \
     -w "\nStatus: %{http_code}\n"
echo "----------------------------------------"
sleep 1

# 6. Удаление заметки
echo ""
echo "Запрос 6: DELETE $NOTES_URL/$NOTE_ID"

curl -X DELETE "$NOTES_URL/$NOTE_ID" \
     -H "Authorization: Bearer $ACCESS_TOKEN" \
     -w "\nStatus: %{http_code}\n"
echo "----------------------------------------"

echo ""
echo "Тестирование завершено!"