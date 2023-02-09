package functions

import (
	"fmt"
)

func SendLeave(user_id string,from_date string,to_date string) {
	fmt.Println("In sendLeave",user_id,from_date,to_date)
}
func SendCancelLeave(user_id string) {
	fmt.Println("In sendCancelLeave",user_id)
}
func GetInsights(user_id string,from_date string,to_date string,format string) string {
	fmt.Println("In getInsights")
	return "Insights for the period "+from_date+" to "+to_date+" in "+format+" format"
}