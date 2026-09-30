AUTH_URL="http://localhost/auth"
NOTES_URL="http://localhost/notes"

echo "=== Подготовка: Регистрация и получение токена через Nginx ==="
TEST_USER="gateway_tester_$RANDOM"

echo ""
echo "Запрос 1: POST $AUTH_URL/register"
curl -X POST "$AUTH_URL/register" \
     -H "Content-Type: application/json" \
     -d "{\"username\": \"$TEST_USER\", \"password\": \"testpass123\"}" \
     -w "\nStatus: %{http_code}\n" \
     -s
echo "----------------------------------------"
sleep 1

echo ""
echo "Запрос 2: POST $AUTH_URL/login"
LOGIN_RESPONSE=$(curl -s -X POST "$AUTH_URL/login" \
     -H "Content-Type: application/json" \
     -d "{\"username\": \"$TEST_USER\", \"password\": \"testpass123\"}")

echo "$LOGIN_RESPONSE"
ACCESS_TOKEN=$(echo "$LOGIN_RESPONSE" | grep -o '"access_token":"[^"]*"' | cut -d'"' -f4)
REFRESH_TOKEN=$(echo "$LOGIN_RESPONSE" | grep -o '"refresh_token":"[^"]*"' | cut -d'"' -f4)

if [ -z "$ACCESS_TOKEN" ]; then
    echo "Ошибка: Не удалось получить токен. Проверь роутинг Nginx."
    exit 1
fi

echo "Успешно получен токен для $TEST_USER"
echo "----------------------------------------"
sleep 1

echo ""
echo "Запрос 3: GET $AUTH_URL/user"
curl -X GET "$AUTH_URL/user" \
     -H "Authorization: Bearer $ACCESS_TOKEN" \
     -w "\nStatus: %{http_code}\n" \
     -s
echo "----------------------------------------"
sleep 1

echo "=== Тестирование Notes API через Nginx ==="
echo "----------------------------------------"

echo ""
echo "Запрос 4: POST $NOTES_URL"
CREATE_RESPONSE=$(curl -X POST "$NOTES_URL" \
     -H "Authorization: Bearer $ACCESS_TOKEN" \
     -H "Content-Type: application/json" \
     -d '{
           "name": "Nginx Gateway Test",
           "content": "Этот запрос прошел через единую точку входа"
         }' \
     -w "\nStatus: %{http_code}\n" \
     -s)

echo "$CREATE_RESPONSE"
NOTE_ID=$(echo "$CREATE_RESPONSE" | grep -o '"id":"[^"]*"' | cut -d'"' -f4)

echo "Извлеченный ID заметки: $NOTE_ID"
echo "----------------------------------------"
sleep 1

echo ""
echo "Запрос 5: GET $NOTES_URL"
curl -X GET "$NOTES_URL" \
     -H "Authorization: Bearer $ACCESS_TOKEN" \
     -w "\nStatus: %{http_code}\n" \
     -s
echo "----------------------------------------"
sleep 1

echo ""
echo "Запрос 6: GET $NOTES_URL/$NOTE_ID"
curl -X GET "$NOTES_URL/$NOTE_ID" \
     -H "Authorization: Bearer $ACCESS_TOKEN" \
     -w "\nStatus: %{http_code}\n" \
     -s
echo "----------------------------------------"
sleep 1

echo ""
echo "Запрос 7: PATCH $NOTES_URL/$NOTE_ID"
curl -X PATCH "$NOTES_URL/$NOTE_ID" \
     -H "Authorization: Bearer $ACCESS_TOKEN" \
     -H "Content-Type: application/json" \
     -d '{
           "name": "Обновлено через Gateway",
           "content": "Проверяем метод PATCH"
         }' \
     -w "\nStatus: %{http_code}\n" \
     -s
echo "----------------------------------------"
sleep 1

echo ""
echo "Запрос 8: DELETE $NOTES_URL/$NOTE_ID"
curl -X DELETE "$NOTES_URL/$NOTE_ID" \
     -H "Authorization: Bearer $ACCESS_TOKEN" \
     -w "\nStatus: %{http_code}\n" \
     -s
echo "----------------------------------------"
sleep 1

echo "=== Завершение: Проверка Refresh и удаление пользователя ==="
echo "----------------------------------------"

echo ""
echo "Запрос 9: POST $AUTH_URL/refresh"
curl -X POST "$AUTH_URL/refresh" \
     -H "Content-Type: application/json" \
     -d "{\"refresh_token\": \"$REFRESH_TOKEN\"}" \
     -w "\nStatus: %{http_code}\n" \
     -s
echo "----------------------------------------"
sleep 1

echo ""
echo "Запрос 10: DELETE $AUTH_URL/user"
curl -X DELETE "$AUTH_URL/user" \
     -H "Authorization: Bearer $ACCESS_TOKEN" \
     -w "\nStatus: %{http_code}\n" \
     -s
echo "----------------------------------------"

echo ""
echo "Глобальное тестирование завершено!"