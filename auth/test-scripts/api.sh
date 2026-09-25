BASE_URL="http://localhost:8081/auth"

echo "Тестирование Auth API"
echo "----------------------------------------"

# 1. Регистрация нового пользователя
echo ""
echo "Запрос 1: POST $BASE_URL/register"

curl -X POST "$BASE_URL/register" \
     -H "Content-Type: application/json" \
     -d '{
           "username": "testuser3",
           "password": "password123"
         }' \
     -w "\nStatus: %{http_code}\n"
echo "----------------------------------------"

sleep 1

# 2. Логин пользователя
echo ""
echo "Запрос 2: POST $BASE_URL/login"

LOGIN_RESPONSE=$(curl -X "POST" "$BASE_URL/login" \
     -H "Content-Type: application/json" \
     -d '{
           "username": "testuser3",
           "password": "password123"
         }' \
     -w "\nStatus: %{http_code}\n" \
     -s)
echo $LOGIN_RESPONSE

ACCESS_TOKEN_1=$(echo "$LOGIN_RESPONSE" | grep -o '"access_token":"[^"]*"' | cut -d'"' -f4)
REFRESH_TOKEN=$(echo "$LOGIN_RESPONSE" | grep -o '"refresh_token":"[^"]*"' | cut -d'"' -f4)
echo "Извлеченный токен доступа: $ACCESS_TOKEN_1"
echo "Извлеченный токен обновления: $REFRESH_TOKEN"
echo "----------------------------------------"

sleep 1

# 3. Получение данных о пользователе
echo ""
echo "Запрос 3: GET $BASE_URL/user"

curl -X GET "$BASE_URL/user" \
     -H "Authorization: Bearer $ACCESS_TOKEN_1" \
     -H "Content-Type: application/json" \
     -w "\nStatus: %{http_code}\n"
echo "----------------------------------------"

sleep 1

# 4. Обновление данных пользователя
echo ""
echo "Запрос 4: PATCH $BASE_URL/user"

curl -X PATCH "$BASE_URL/user" \
     -H "Content-Type: application/json" \
     -H "Authorization: Bearer $ACCESS_TOKEN_1" \
     -d '{
          "password": "updated@example.com"
         }' \
     -w "\nStatus: %{http_code}\n"
echo "----------------------------------------"

sleep 1

# 5. Проверка обновленной информации
echo ""
echo "Запрос 5: GET $BASE_URL/user"

curl -X GET "$BASE_URL/user" \
     -H "Authorization: Bearer $ACCESS_TOKEN_1" \
     -H "Content-Type: application/json" \
     -w "\nStatus: %{http_code}\n"
echo "----------------------------------------"

sleep 1

# 6. Рефреш Access Токена
echo ""
echo "Запрос 6: POST $BASE_URL/refresh"

REFRESH_RESPONSE=$(curl -X "POST" "$BASE_URL/refresh" \
     -H "Content-Type: application/json" \
     -d "{\"refresh_token\": \"$REFRESH_TOKEN\"}" \
     -w "\nStatus: %{http_code}\n" \
     -s)
echo $REFRESH_RESPONSE

ACCESS_TOKEN_2=$(echo "$REFRESH_RESPONSE" | grep -o '"access_token":"[^"]*"' | cut -d'"' -f4)
echo "Извлеченный токен доступа: $ACCESS_TOKEN_2"
echo "----------------------------------------"

sleep 1

# 7. Удаление пользователя
echo ""
echo "Запрос 7: DELETE $BASE_URL/user"

curl -X DELETE "$BASE_URL/user" \
     -H "Authorization: Bearer $ACCESS_TOKEN_2" \
     -H "Content-Type: application/json" \
     -w "\nStatus: %{http_code}\n"
echo "----------------------------------------"

echo ""
echo "Тестирование завершено!"