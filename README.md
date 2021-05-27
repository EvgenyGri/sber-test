#TEST EXAM

##usage
```
sber-test --source "file-name.json" 
          [--count-per-recipe]
          [--unique-recipe-count]
          [--busiest-postcode]
          [--find-recipes  "name1,name1,.."]
          [--deliveries-by-postcode-and-time "postcode,from,to"]
```

##example
```
sber-test --source ./data.json\
          --count-per-recipe\ 
          --busiest-postcode\
          --unique-recipe-count\
          --find-recipes "Chicken,Cherry,Tilapia"\
          --deliveries-by-postcode-and-time "10163,6AM,6PM"
```
##параметры:
- ```--source```  указывает на файл
- ```--count-per-recipe``` Подсчитать число вхождений каждого уникального "recipe name" (с алфавитной сортировкой по "recipe name")
- ```--busiest-postcode``` Подсчитать число уникальных "recipe name"
- ```--busiest-postcode``` Найти "postcode" с наибольшим числом доаставок.
- ```--find-recipes``` Перечислить "recipe name" (в алфавитном порядке), которые содержат в своём имени одно из слов
- ```--deliveries-by-postcode-and-time``` Найти число доставок для "postcode", которые происходили во временном промежутке "from.to" 
            