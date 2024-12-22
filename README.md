Начало от localhost:port/quiz

Добавление/Изменение вариантов/вопросов осуществляется по Postman

**Если сервис запускается в сети docker-compose, то в config.json db.host должен быть равен названию контейнера постгреса. postgres - по дефолту, localhost - иначе**

**Заполнить .env при необходимости**

- [ GET ]    -->      /quiz/                    
- [ POST ]   -->      /quiz/register
```
Body:
{
    Login    string `json:"login" binding:"required"`
    Password string `json:"password" binding:"required"`
}
```

- [ POST ]   -->      /quiz/login
```
Body:
{
    Login    string `json:"login" binding:"required"`
    Password string `json:"password" binding:"required"`
}
```
- [ POST ]   -->      /quiz/:userId/quit        
- [ GET ]    -->      /quiz/:userId/variant/    
- [ POST ]   -->      /quiz/:userId/variant/add
```
Body:
{
    Name    string `json:"name" binding:"required"`
}
```
- [ GET ]    -->      /quiz/:userId/variant/list 
- [ GET ]    -->      /quiz/:userId/variant/:variantName/
- [ DELETE ] -->      /quiz/:userId/variant/:variantName/remove 
- [ POST ]   -->      /quiz/:userId/variant/:variantName/start 
- [ GET ]    -->      /quiz/:userId/variant/:variantName/results 
- [ GET ]    -->      /quiz/:userId/variant/:variantName/get 
- [ POST ]   -->      /quiz/:userId/variant/:variantName/question/add 
```
Body:
{
	Question string    `json:"question" binding:"required,max=50" db:"question"`
	Answer   string    `json:"answer" binding:"required,max=50" db:"answer"`
	Answers  []*Answer `json:"answers" binding:"required,len=3"`
}
Example:
{
    "question": "Первый вопрос",
    "answer": "правильный ответ",
    "answers": [
        {
            "answer": "неправильный ответ 1"
        },
        {
            "answer": "неправильный ответ 2"
        },
        {
            "answer": "неправильный ответ 3"
        }
    ]
}
```
- [ DELETE ] -->      /quiz/:userId/variant/:variantName/question/remove 
- [ GET ]    -->      /quiz/:userId/variant/:variantName/question/:questionId/get 
- [ POST ]   -->      /quiz/:userId/variant/:variantName/question/:questionId/accept 
```
Body:
{
    Answer string `json:"answer" binding:"required"`
}
```