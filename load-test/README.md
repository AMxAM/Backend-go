# Configuración de Pruebas de Carga

Herramienta:
- K6 v2.0.0

Endpoint Evaluado:
- POST /api/v1/notifications/send

Infraestructura:
- API Gateway
- AWS Lambda (go-api-hex)
- SNS
- SQS
- Notification Lambda
- SMTP Gmail

Escenarios:

1. Línea Base
   - 10 usuarios
   - 1 minuto

2. Incremento Progresivo
   - 50 usuarios
   - 2 minutos

3. Incremento Progresivo
   - 100 usuarios
   - 2 minutos

4. Incremento Progresivo
   - 250 usuarios
   - 2 minutos

5. Incremento Progresivo
   - 500 usuarios
   - 2 minutos