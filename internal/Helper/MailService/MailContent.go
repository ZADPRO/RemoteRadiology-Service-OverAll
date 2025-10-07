package mailservice

import (
	timeZone "AuthenticationService/internal/Helper/TimeZone"
	"html"
	"strconv"
	"time"
)

func LoginOTPContent(otp int) string {
	return `
	<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <title>Login Verification Code</title>
  <style>
    body {
      font-family: Arial, sans-serif;
      background-color: #f6f8fa;
      margin: 0;
      padding: 0;
    }
    .container {
      max-width: 600px;
      background-color: #EDD1CE;
      margin: 40px auto;
      padding: 30px;
      border-radius: 8px;
      box-shadow: 0 0 10px rgba(0,0,0,0.05);
    }
    .header {
      text-align: center;
      padding-bottom: 20px;
    }
    .header h1 {
      margin: 0;
      color: #525252;
    }
    .otp {
      font-size: 32px;
      font-weight: 700;
      color: #333333;
      text-align: center;
      margin: 20px 0;
      letter-spacing: 2px;
    }
    .content {
      font-size: 16px;
      color: #525252;
      text-align: center;
      margin-bottom: 30px;
    }
    .footer {
      font-size: 12px;
      text-align: center;
      color: #525252;
      border-top: 1px solid #dddddd;
      padding-top: 15px;
    }
  </style>
</head>
<body>
  <div class="container">
    <div class="header">
      <h1>Wellthgreen Report Portal Login Verification Code</h1>
    </div>
    <div class="content">
      <p>Please use the Wellthgreen Report Portal verification code below to log in. This code will remain valid for 10 minutes:</p>
      <div class="otp">` + html.EscapeString(strconv.Itoa(otp)) + `</div>
      <p style="text-align: justify;">This email is intended only for the individual or entity to which it is addressed, and may contain information that is privileged, confidential, and exempt from disclosure under applicable law. If the reader of this message is not the intended recipient, or the employee or agent responsible for delivering the message to the intended recipient, you are hereby informed that any use, disclosure, distribution, or copying of this communication is strictly prohibited. If you have received this communication in error, please notify us immediately by telephone and delete the original email message.</p>
    </div>
    <div class="footer">
      &copy; ` + html.EscapeString(strconv.Itoa(time.Now().In(timeZone.MustGetPacificLocation()).Year())) + ` Wellthgreen. All rights reserved.
    </div>
  </div>
</body>
</html>
	`
}

func ForgetPasswordOTPContent(otp int) string {
	return `
			<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <title>Login Verification Code</title>
  <style>
    body {
      font-family: Arial, sans-serif;
      background-color: #f6f8fa;
      margin: 0;
      padding: 0;
    }
    .container {
      max-width: 600px;
      background-color: #EDD1CE;
      margin: 40px auto;
      padding: 30px;
      border-radius: 8px;
      box-shadow: 0 0 10px rgba(0,0,0,0.05);
    }
    .header {
      text-align: center;
      padding-bottom: 20px;
    }
    .header h1 {
      margin: 0;
      color: #525252;
    }
    .otp {
      font-size: 32px;
      font-weight: bold;
      color: #333333;
      text-align: center;
      margin: 20px 0;
      letter-spacing: 2px;
    }
    .content {
      font-size: 16px;
      color: #525252;
      text-align: center;
      margin-bottom: 30px;
    }
    .footer {
      font-size: 12px;
      text-align: center;
      color: #525252;
      border-top: 1px solid #dddddd;
      padding-top: 15px;
    }
  </style>
</head>
<body>
  <div class="container">
    <div class="header">
      <h1>Wellthgreen Report Portal Password Reset Code</h1>
    </div>
    <div class="content">
      <p>Please use the Wellthgreen Report Portal verification code below to reset your password. This code will remain valid for 10 minutes:</p>
      <div class="otp">` + html.EscapeString(strconv.Itoa(otp)) + `</div>
      <p style="text-align: justify;">This email is intended only for the individual or entity to which it is addressed, and may contain information that is privileged, confidential, and exempt from disclosure under applicable law. If the reader of this message is not the intended recipient, or the employee or agent responsible for delivering the message to the intended recipient, you are hereby informed that any use, disclosure, distribution, or copying of this communication is strictly prohibited. If you have received this communication in error, please notify us immediately by telephone and delete the original email message.</p>
    </div>
    <div class="footer">
      &copy; ` + html.EscapeString(strconv.Itoa(time.Now().In(timeZone.MustGetPacificLocation()).Year())) + ` Wellthgreen. All rights reserved.
    </div>
  </div>
</body>
</html>
	`
}

func RegisterationMailContent(userName string) string {
	return `
<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="UTF-8" />
    <style>
      body {
        font-family: Arial, sans-serif;
        background-color: #f6f8fa;
        margin: 0;
        padding: 0;
      }
      .container {
        max-width: 600px;
        background-color: #EDD1CE;
        margin: 40px auto;
        padding: 30px;
        border-radius: 8px;
        box-shadow: 0 0 10px rgba(0, 0, 0, 0.05);
      }
      .header {
        text-align: center;
        padding-bottom: 20px;
      }
      .header h1 {
        margin: 0;
        color: #525252;
      }
      .content {
        font-size: 16px;
        color: #525252;
        text-align: center;
        margin-bottom: 30px;
      }
      .button {
        display: inline-block;
        padding: 12px 25px;
        background-color: #c6d4c0;
        color: #f7f7f7;
        border-radius: 5px;
        text-decoration: none;
        font-weight: bold;
        margin-top: 20px;
      }
      .footer {
        font-size: 12px;
        text-align: center;
        color: #525252;
        border-top: 1px solid #dddddd;
        padding-top: 15px;
      }
    </style>
  </head>
  <body>
    <div class="container">
      <div class="header">
        <h1>Welcome, ` + html.EscapeString(userName) + `!</h1>
      </div>
      <div class="content">
        <p>
          Thank you for registering with <strong>Wellthgreen</strong>.
        </p>
        <p>
          We’re excited to have you on board. You can now log in using your
          registered email and the password you created.
        </p>
        <a href="https://reportportal.wellthgreen.com/" class="button">Login Now</a>
        <p style="text-align: justify;">This email is intended only for the individual or entity to which it is addressed, and may contain information that is privileged, confidential, and exempt from disclosure under applicable law. If the reader of this message is not the intended recipient, or the employee or agent responsible for delivering the message to the intended recipient, you are hereby informed that any use, disclosure, distribution, or copying of this communication is strictly prohibited. If you have received this communication in error, please notify us immediately by telephone and delete the original email message.</p>
      </div>
      <div class="footer">
        &copy; ` + html.EscapeString(strconv.Itoa(time.Now().In(timeZone.MustGetPacificLocation()).Year())) + `
        Wellthgreen. All rights reserved.
      </div>
    </div>
  </body>
</html>
  `
}

func GetOTPMailContent(userName string, otp int) string {
	return `
 <!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="UTF-8" />
    <style>
      body {
        font-family: Arial, sans-serif;
        background-color: #f6f8fa;
        margin: 0;
        padding: 0;
      }
      .container {
        max-width: 600px;
        background-color: #edd1ce;
        margin: 40px auto;
        padding: 30px;
        border-radius: 8px;
        box-shadow: 0 0 10px rgba(0, 0, 0, 0.05);
      }
      .header {
        text-align: center;
        padding-bottom: 20px;
      }
      .header h1 {
        margin: 0;
        color: #525252;
      }
      .otp-box {
        font-size: 24px;
        font-weight: bold;
        background-color: #f0f0f0;
        padding: 15px;
        border-radius: 8px;
        margin: 20px auto;
        text-align: center;
        color: #333333;
        letter-spacing: 4px;
        max-width: 200px;
      }
      .content {
        font-size: 16px;
        color: #525252;
        text-align: center;
        margin-bottom: 30px;
      }
      .button {
        display: inline-block;
        padding: 12px 25px;
        background-color: #c6d4c0;
        color: #ffffff;
        border-radius: 5px;
        text-decoration: none;
        font-weight: bold;
        margin-top: 20px;
      }
      .footer {
        font-size: 12px;
        text-align: center;
        color: #525252;
        border-top: 1px solid #dddddd;
        padding-top: 15px;
      }
    </style>
  </head>
  <body>
    <div class="container">
      <div class="header">
        <h1>Welcome, ` + html.EscapeString(userName) + `!</h1>
      </div>
      <div class="content">
        <p>
          Thank you for registering with
          <strong>Wellthgreen Report Portal</strong>.
        </p>
        <p>
          Please use the Verification Code below to verify your email address and complete
          your registration:
        </p>
        <div class="otp-box">` + html.EscapeString(strconv.Itoa(otp)) + `</div>
        <p>This Verification Code is valid for the next 10 minutes.</p>
        <p style="text-align: justify;">This email is intended only for the individual or entity to which it is addressed, and may contain information that is privileged, confidential, and exempt from disclosure under applicable law. If the reader of this message is not the intended recipient, or the employee or agent responsible for delivering the message to the intended recipient, you are hereby informed that any use, disclosure, distribution, or copying of this communication is strictly prohibited. If you have received this communication in error, please notify us immediately by telephone and delete the original email message.</p>
      </div>
      <div class="footer">
        &copy; ` + html.EscapeString(strconv.Itoa(time.Now().In(timeZone.MustGetPacificLocation()).Year())) + `
        Wellthgreen. All rights reserved.
      </div>
    </div>
  </body>
</html>
  `
}

func RegistrationMailContent(userName, patientID, gmail, password string, role string) string {
	return `
<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="UTF-8" />
    <style>
      body {
        font-family: Arial, sans-serif;
        background-color: #f6f8fa;
        margin: 0;
        padding: 0;
      }
      .container {
        max-width: 600px;
        background-color: #edd1ce;
        margin: 40px auto;
        padding: 30px;
        border-radius: 8px;
        box-shadow: 0 0 10px rgba(0, 0, 0, 0.05);
      }
      .header {
        text-align: center;
        padding-bottom: 20px;
      }
      .header h1 {
        margin: 0;
        color: #525252;
      }
      .content {
        font-size: 16px;
        color: #525252;
        text-align: center;
        margin-bottom: 30px;
      }
      .credentials {
        background-color: #fff;
        padding: 15px;
        border-radius: 5px;
        margin: 20px auto;
        width: fit-content;
        text-align: left;
        font-family: monospace;
        border: 1px solid #ccc;
      }
      .button {
        display: inline-block;
        padding: 12px 25px;
        background-color: #c6d4c0;
        color: #ffffff;
        border-radius: 5px;
        text-decoration: none;
        font-weight: bold;
        margin-top: 20px;
      }
      .footer {
        font-size: 12px;
        text-align: center;
        color: #525252;
        border-top: 1px solid #dddddd;
        padding-top: 15px;
      }
    </style>
  </head>
  <body>
    <div class="container">
      <div class="header">
        <h1>Welcome, ` + html.EscapeString(userName) + `!</h1>
      </div>
      <div class="content">
        <p>
          You have successfully been onboarded as a ` + html.EscapeString(role) + ` on <strong>Wellthgreen Report Portal</strong>.
        </p>
        <p>Your login credentials are as follows:</p>
        <div class="credentials">
          <p>
            <strong>` + html.EscapeString(role) + ` ID:</strong> ` +
		html.EscapeString(patientID) + `
          </p>
          <p>
            <strong>Email:</strong> ` + html.EscapeString(gmail) + `
          </p>
          <p><strong>Password:</strong> ` + html.EscapeString(password) + `</p>
        </div>
        <a href="https://reportportal.wellthgreen.com/" class="button"
          >Login Now</a
        >
        <p style="text-align: justify;margin-top: 20px">This email is intended only for the individual or entity to which it is addressed, and may contain information that is privileged, confidential, and exempt from disclosure under applicable law. If the reader of this message is not the intended recipient, or the employee or agent responsible for delivering the message to the intended recipient, you are hereby informed that any use, disclosure, distribution, or copying of this communication is strictly prohibited. If you have received this communication in error, please notify us immediately by telephone and delete the original email message.</p>
      </div>
      <div class="footer">
        &copy; ` + html.EscapeString(strconv.Itoa(time.Now().In(timeZone.MustGetPacificLocation()).Year())) + `
        Wellthgreen. All rights reserved.
      </div>
    </div>
  </body>
</html>

`
}

func PatientReportSignOff(userName string, patientID string, AppintmentDate string, scancenterCode string) string {
	return `
  <!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="UTF-8" />
    <style>
      body {
        font-family: Arial, sans-serif;
        background-color: #f6f8fa;
        margin: 0;
        padding: 0;
      }
      .container {
        max-width: 600px;
        background-color: #edd1ce;
        margin: 40px auto;
        padding: 30px;
        border-radius: 8px;
        box-shadow: 0 0 10px rgba(0, 0, 0, 0.05);
      }
      .header {
        text-align: center;
        padding-bottom: 20px;
      }
      .header h1 {
        margin: 0;
        color: #525252;
      }
      .content {
        font-size: 16px;
        color: #525252;
        text-align: center;
        margin-bottom: 30px;
      }
      .report-info {
        background-color: #fff;
        padding: 15px;
        border-radius: 5px;
        margin: 20px auto;
        width: fit-content;
        text-align: left;
        font-family: monospace;
        border: 1px solid #ccc;
      }
      .button {
        display: inline-block;
        padding: 12px 25px;
        background-color: #c6d4c0;
        color: #ffffff;
        border-radius: 5px;
        text-decoration: none;
        font-weight: bold;
        margin-top: 20px;
      }
      .highlight {
        background-color: #e8f5e8;
        padding: 15px;
        border-radius: 5px;
        border-left: 4px solid #28a745;
        margin: 20px 0;
        font-weight: bold;
        color: #155724;
      }
      .footer {
        font-size: 12px;
        text-align: center;
        color: #525252;
        border-top: 1px solid #dddddd;
        padding-top: 15px;
      }
    </style>
  </head>
  <body>
    <div class="container">
      <div class="header">
        <h1>Report Ready - ` + html.EscapeString(userName) + `!</h1>
      </div>
      <div class="content">
        <div class="highlight">
          📋 Your scan report has been completed!
        </div>
        <p>
          Dear ` + html.EscapeString(userName) + `,
        </p>
        <p>
          Your report has been processed and is now available for download.
        </p>
        
        <div class="report-info">
          <p><strong>Patient ID:</strong> ` + html.EscapeString(patientID) + `</p>
          <p><strong>Appointment Date:</strong> ` + html.EscapeString(AppintmentDate) + `</p>
          <p><strong>Scan Center Code:</strong> ` + html.EscapeString(scancenterCode) + `</p>
        </div>

        <p>To access your report, please log in with your credentials:</p>
        
        <a href="https://reportportal.wellthgreen.com/" class="button">
          Login to View Report
        </a>
        
        <p style="margin-top: 15px; font-size: 14px;text-align:justify;">
          This email is intended only for the individual or entity to which it is addressed, and may contain information that is privileged, confidential, and exempt from disclosure under applicable law. If the reader of this message is not the intended recipient, or the employee or agent responsible for delivering the message to the intended recipient, you are hereby informed that any use, disclosure, distribution, or copying of this communication is strictly prohibited. If you have received this communication in error, please notify us immediately by telephone and delete the original email message.
        </p>
      </div>
      <div class="footer">
        &copy; ` + html.EscapeString(strconv.Itoa(time.Now().In(timeZone.MustGetPacificLocation()).Year())) + `
        Wellthgreen. All rights reserved.<br>
      </div>
    </div>
  </body>
</html>
  `
}

func ManagerReportSignOff(patientName string, patientID string, appointmentDate string, scanCenterCode string) string {
	return `
  <!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="UTF-8" />
    <style>
      body {
        font-family: Arial, sans-serif;
        background-color: #f6f8fa;
        margin: 0;
        padding: 0;
      }
      .container {
        max-width: 600px;
        background-color: #edd1ce;
        margin: 40px auto;
        padding: 30px;
        border-radius: 8px;
        box-shadow: 0 0 10px rgba(0, 0, 0, 0.05);
      }
      .header {
        text-align: center;
        padding-bottom: 20px;
      }
      .header h1 {
        margin: 0;
        color: #525252;
      }
      .content {
        font-size: 16px;
        color: #525252;
        text-align: center;
        margin-bottom: 30px;
      }
      .report-info {
        background-color: #fff;
        padding: 15px;
        border-radius: 5px;
        margin: 20px auto;
        width: fit-content;
        text-align: left;
        font-family: monospace;
        border: 1px solid #ccc;
      }
      .button {
        display: inline-block;
        padding: 12px 25px;
        background-color: #c6d4c0;
        color: #ffffff;
        border-radius: 5px;
        text-decoration: none;
        font-weight: bold;
        margin-top: 20px;
      }
      .highlight {
        background-color: #e8f5e8;
        padding: 15px;
        border-radius: 5px;
        border-left: 4px solid #28a745;
        margin: 20px 0;
        font-weight: bold;
        color: #155724;
      }
      .footer {
        font-size: 12px;
        text-align: center;
        color: #525252;
        border-top: 1px solid #dddddd;
        padding-top: 15px;
      }
    </style>
  </head>
  <body>
    <div class="container">
      <div class="header">
        <h1>Report Ready - ` + html.EscapeString(patientName) + `!</h1>
      </div>
      <div class="content">
        <div class="highlight">
          📋 Scan report has been completed!
        </div>
        <p>
          Report has been processed and is now available for download.
        </p>
        
        <div class="report-info">
          <p><strong>Patient ID:</strong> ` + html.EscapeString(patientID) + `</p>
          <p><strong>Appointment Date:</strong> ` + html.EscapeString(appointmentDate) + `</p>
          <p><strong>Scan Center Code:</strong> ` + html.EscapeString(scanCenterCode) + `</p>
        </div>

        <p>To access your report, please log in with your credentials:</p>
        
        <a href="https://reportportal.wellthgreen.com/" class="button">
          Login to View Report
        </a>
        
        <p style="margin-top: 15px; font-size: 14px;text-align:justify;">
          This email is intended only for the individual or entity to which it is addressed, and may contain information that is privileged, confidential, and exempt from disclosure under applicable law. If the reader of this message is not the intended recipient, or the employee or agent responsible for delivering the message to the intended recipient, you are hereby informed that any use, disclosure, distribution, or copying of this communication is strictly prohibited. If you have received this communication in error, please notify us immediately by telephone and delete the original email message.
        </p>
      </div>
      <div class="footer">
        &copy; ` + html.EscapeString(strconv.Itoa(time.Now().In(timeZone.MustGetPacificLocation()).Year())) + `
        Wellthgreen. All rights reserved.<br>
      </div>
    </div>
  </body>
</html>
 `
}
