AUTH_URL="http://localhost:8081/auth"
NOTES_URL="http://localhost:8082/notes"
FAKE_ID="60d5ec49c2a23e5a5a123456" 

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
    echo "Ошибка: Не удалось получить токен."
    exit 1
fi

echo "Успешно получен токен для $TEST_USER"
echo "----------------------------------------"
sleep 1

echo "=== Успешные сценарии (200 / 201) ==="

echo -e "\nЗапрос 1: POST $NOTES_URL (Создание заметки)"
CREATE_RESPONSE=$(curl -X POST "$NOTES_URL" \
     -H "Authorization: Bearer $ACCESS_TOKEN" \
     -H "Content-Type: application/json" \
     -d '{"name": "Финальный тест", "content": "Всё работает идеально"}' \
     -w "\nStatus: %{http_code}\n" -s)
echo "$CREATE_RESPONSE"
NOTE_ID=$(echo "$CREATE_RESPONSE" | grep -o '"id":"[^"]*"' | cut -d'"' -f4)

echo -e "\nЗапрос 2: GET $NOTES_URL (Получение всех заметок)"
curl -X GET "$NOTES_URL" \
     -H "Authorization: Bearer $ACCESS_TOKEN" \
     -w "\nStatus: %{http_code}\n" -s

echo -e "\nЗапрос 3: GET $NOTES_URL/$NOTE_ID (Получение по ID)"
curl -X GET "$NOTES_URL/$NOTE_ID" \
     -H "Authorization: Bearer $ACCESS_TOKEN" \
     -w "\nStatus: %{http_code}\n" -s

echo -e "\nЗапрос 4: PATCH $NOTES_URL/$NOTE_ID (Частичное обновление - только контент)"
curl -X PATCH "$NOTES_URL/$NOTE_ID" \
     -H "Authorization: Bearer $ACCESS_TOKEN" \
     -H "Content-Type: application/json" \
     -d '{"content": "Обновили только текст, заголовок остался"}' \
     -w "\nStatus: %{http_code}\n" -s

echo "----------------------------------------"
echo "=== Тестирование ошибок (400 / 404) ==="

echo -e "\nЗапрос 5: POST $NOTES_URL (Ошибка 400 - Пустое имя)"
curl -X POST "$NOTES_URL" \
     -H "Authorization: Bearer $ACCESS_TOKEN" \
     -H "Content-Type: application/json" \
     -d '{"content": "Заметка без имени"}' \
     -w "\nStatus: %{http_code}\n" -s

echo -e "\nЗапрос 6: PATCH $NOTES_URL/$FAKE_ID (Ошибка 404 - Обновление несуществующей)"
curl -X PATCH "$NOTES_URL/$FAKE_ID" \
     -H "Authorization: Bearer $ACCESS_TOKEN" \
     -H "Content-Type: application/json" \
     -d '{"name": "Призрак"}' \
     -w "\nStatus: %{http_code}\n" -s

echo -e "\nЗапрос 7: GET $NOTES_URL/$FAKE_ID (Ошибка 404 - Получение несуществующей)"
curl -X GET "$NOTES_URL/$FAKE_ID" \
     -H "Authorization: Bearer $ACCESS_TOKEN" \
     -w "\nStatus: %{http_code}\n" -s

echo -e "\nЗапрос 8: GET $NOTES_URL/invalid-id-format (Ошибка 400 - Кривой ID)"
curl -X GET "$NOTES_URL/invalid-id-format" \
     -H "Authorization: Bearer $ACCESS_TOKEN" \
     -w "\nStatus: %{http_code}\n" -s

echo "----------------------------------------"
echo "=== Очистка данных ==="

echo -e "\nЗапрос 9: DELETE $NOTES_URL/$NOTE_ID (Успешное удаление)"
curl -X DELETE "$NOTES_URL/$NOTE_ID" \
     -H "Authorization: Bearer $ACCESS_TOKEN" \
     -w "\nStatus: %{http_code}\n" -s

echo -e "\nЗапрос 10: DELETE $NOTES_URL/$NOTE_ID (Ошибка 404 - Повторное удаление)"
curl -X DELETE "$NOTES_URL/$NOTE_ID" \
     -H "Authorization: Bearer $ACCESS_TOKEN" \
     -w "\nStatus: %{http_code}\n" -s

echo -e "\nТестирование завершено!"