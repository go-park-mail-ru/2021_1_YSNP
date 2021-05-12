FROM build

FROM alpine:latest

COPY --from=build /app/main /app/

WORKDIR /app

CMD ./main