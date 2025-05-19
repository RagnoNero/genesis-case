package email

import (
	"fmt"
	"weather-subscription/models"
)

func generateConfirmEmailBody(baseUrl, token string) string {
	return fmt.Sprintf(`
		<html>
		<body style="font-family: Arial, sans-serif; line-height: 1.6;">
			<h2>Please confirm your subscription</h2>
			<p>Click the button below to confirm your subscription:</p>
			<a href="%s/confirm/%s" 
				style="
					display: inline-block;
					padding: 12px 24px;
					font-size: 16px;
					color: white;
					background-color: #007BFF;
					text-decoration: none;
					border-radius: 6px;
				"
			>Confirm</a>
			<p>If you did not request this subscription, you can ignore this email.</p>
			<p style="font-size: 14px; color: #555;">To unsubscribe from future emails, click <a href="%s/unsubscribe/%s">here</a>.</p>
		</body>
		</html>
	`, baseUrl, token, baseUrl, token)
}

func generateWeatherEmailBody(baseUrl, city string, weather models.Weather, token string) string {
	return fmt.Sprintf(`
		<html>
		<body style="font-family: Arial, sans-serif; line-height: 1.6;">
			<h2>Weather update for %s</h2>
			<p><strong>Temperature:</strong> %d°C</p>
			<p><strong>Humidity:</strong> %d%%</p>
			<p><strong>Conditions:</strong> %s</p>

			<hr style="margin: 30px 0;">

			<p style="font-size: 14px; color: #555;">
				If you no longer wish to receive weather updates, you can 
				<a href="%s/unsubscribe/%s" style="color: #d00;">unsubscribe here</a>.
			</p>
		</body>
		</html>
	`, city, weather.Temperature, weather.Humidity, weather.Description, baseUrl, token)
}
