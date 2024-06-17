# go-tracing

Para criar a imagem no docker: docker build -t go_tracing .

Executar imagem no docker: docker-compose  up -d

Executar os requests, como exemplo na pasta:

- api/service-a-get.http
- api/service-b-post.http

Para acessar o zipkins: http://localhost:9411/zipkin/