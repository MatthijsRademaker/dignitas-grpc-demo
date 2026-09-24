# Build context is the repository root, so the build can see ../proto.
FROM mcr.microsoft.com/dotnet/sdk:10.0 AS build
WORKDIR /src
COPY bff/Bff.csproj bff/
RUN dotnet restore bff/Bff.csproj
COPY proto/ proto/
COPY bff/ bff/
RUN dotnet publish bff/Bff.csproj -c Release -o /app --no-restore

FROM mcr.microsoft.com/dotnet/aspnet:10.0
WORKDIR /app
COPY --from=build /app .
EXPOSE 8080
ENTRYPOINT ["dotnet", "Bff.dll"]
