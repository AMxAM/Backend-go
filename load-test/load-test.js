import http from 'k6/http';

export const options = {
  vus: 500,
  duration: '2m',
};

export default function () {

  const payload = JSON.stringify({
    email: 'test@test.com',
    subject: 'Prueba',
    message: 'Mensaje de carga',
  });

  http.post(
    'https://8g3d9m36sk.execute-api.us-east-2.amazonaws.com/api/v1/notifications/send',
    payload,
    {
      headers: {
        'Content-Type': 'application/json',
      },
    }
  );
}