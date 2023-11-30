package constant

const (
	SmtpAuthAddress   = "smtp.gmail.com"
	SmtpServerAddress = "smtp.gmail.com:587"

	ForgotPwSubject         = "Forgot Password OrenLite"
	ForgotPwContentTemplate = `
	<h3>Hi, %s</h3>
	</p>We've received your request to reset your OrenLite password.</p>
	<br>
	<p>Click <a href="%s">Reset Password</a> to set a new password for your account</p>
	<br>
	<br>
	<p>Have a nice day,<br>OrenLite Team</p>
	`
	ResetPasswordLinkTemplate = "%s/reset-password?code=%s"

	ChangePwSubject         = "Change Password OrenLite"
	ChangePwContentTemplate = `
	<div style="display:flex; flex-direction:column; align-items:center; justify-content:center; text-align:center">
	<h3>Hi, %s</h3>
	<p>We've received your request to cahnge your OrenLite password.
	<br>This is the Verification Code to change you password:</p>
	<br>
	<h1>%s</h1>
	<br>
	<br>
	<p>Have a nice day,<br>OrenLite Team</p>
	</div>
	`
)
