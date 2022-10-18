export const environment = {
  production: true,
  gateway: 'http://' + window.location.hostname,
  cognito: {
    userPoolId: 'us-east-1_wXNyjVJP4',
    userPoolWebClientId: '3bjl2vgb74537vk0b6r33r7jp4',
  },
  mqtt: {
	server: window.location.hostname,
	protocol: "ws",
	port: 9884,
    username: '',
    password: '',
	}
};
