FROM build

FROM alpine:latest

COPY --from=build /app/auth /app/

WORKDIR /app

CMD ./auth