package step

import (
	"context"
	"fmt"
	"visa/global"

	"github.com/chromedp/chromedp"
)

func Login(ctx context.Context) (err error) {
	err = chromedp.Run(ctx,
		chromedp.Navigate("https://ais.usvisa-info.com/en-ca/niv/users/sign_in"),
		chromedp.SendKeys("#user_email", global.CONFIG.Visa.Email, chromedp.ByID),
		chromedp.SendKeys("#user_password", global.CONFIG.Visa.Password, chromedp.ByID),
		chromedp.Click("#policy_confirmed", chromedp.ByID),
		chromedp.Click("[name='commit']", chromedp.ByQuery),
	)
	return err
}

func Step2(ctx context.Context) (err error) {
	err = chromedp.Run(ctx,
		chromedp.WaitVisible("[class='button primary small']", chromedp.ByQuery),
	)
	return err
}

func Step3(ctx context.Context) (err error) {
	err = chromedp.Run(ctx,
		chromedp.Click("[class='button primary small']", chromedp.ByQuery),
	)
	return err
}

func Step4(ctx context.Context) (err error) {
	err = chromedp.Run(ctx,
		chromedp.WaitVisible("[class='accordion-title']", chromedp.ByQuery),
		chromedp.Click("[class='accordion-title']", chromedp.ByQuery),
		chromedp.WaitVisible("[class='button small primary small-only-expanded']", chromedp.ByQuery),
		chromedp.Click("[class='button small primary small-only-expanded']", chromedp.ByQuery),
	)
	return err
}

func Step5(ctx context.Context) (err error) {
	err = chromedp.Run(ctx,
		chromedp.SetValue(`//select[@id="appointments_consulate_appointment_facility_id"]`, "94", chromedp.BySearch),
		chromedp.SetValue(`//select[@id="appointments_consulate_appointment_facility_id"]`, "95", chromedp.BySearch),
	)
	return err
}

func Step6(ctx context.Context, args ...interface{}) (err error) {
	var res any
	targetDate := args[0].(string)

	js := fmt.Sprintf(`
		// Locate the input element by its ID
		let inputElement = document.getElementById("appointments_consulate_appointment_date");

		// Remove the "readonly" attribute
		inputElement.removeAttribute("readonly");

		// Assign a value to the input element
		inputElement.value = "%s";
		setTimeout(function() {
			// click it
			inputElement.click();
			console.log("Start sleeping...");

			// Sleep for 5 seconds (5000 milliseconds)
			setTimeout(function() {
					let dateElement = document.getElementsByClassName("ui-state-default ui-state-active")[0]
					console.log('----var---',dateElement)
					dateElement.click()
			}, 5000);
		}, 1000);



	`, targetDate)
	err = chromedp.Run(ctx,
		chromedp.EvaluateAsDevTools(js, &res),
	)
	return err
}

func Step7(ctx context.Context, args ...interface{}) (err error) {
	var res any
	targetTime := args[0].(string)
	js := fmt.Sprintf(`
	// append options to select
	var selectElement = document.getElementById("appointments_consulate_appointment_time");
	var optionElement = document.createElement("option");
	optionElement.value = "%s";
	optionElement.textContent = "%s";
	selectElement.appendChild(optionElement);

	// Set the value of the time selector
	selectElement.value = "%s";

	// remove disabled button for submit
	document.getElementById("appointments_submit").removeAttribute("disabled");

	// click submit
	document.getElementById("appointments_submit").click();
	`, targetTime, targetTime, targetTime)
	err = chromedp.Run(ctx,
		chromedp.EvaluateAsDevTools(js, &res),
	)
	return err
}
