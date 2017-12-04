package models

type UserReportResult struct {
	Uid          int
	ReportResult ReportResult
}
type ReportResult struct {
	Result    ReResult
	Data      ReData
	Signature string
}
type ReResult struct {
	Success bool
	Message string
}
type KVItem struct {
	Value string
	Key   string
}
type PeerNumTopList struct {
	Top_item []TopItem
	Key      string
}
type FriendCircle struct {
	Summary           []KVItem
	Peer_num_top_list []PeerNumTopList
	Location_top_list []PeerNumTopList
}

type CallRiskAnalysisItem struct {
	Analysis_item  string
	Analysis_point AnalysisPoint
	Analysis_desc  string
}
type AnalysisPoint struct {
	Avg_call_time_6m string
	Call_cnt_1m      string
	Call_cnt_3m      string
	Call_cnt_6m      string
	Call_time_3m     string
	Avg_call_cnt_6m  string
	Call_time_6m     string
	Avg_call_cnt_3m  string
	Call_time_1m     string
	Avg_call_time_3m string
}

type TopItem struct {
	Peer_number     string
	Dialed_time     string
	Call_cnt        string
	Dial_cnt        string
	Dial_time       string
	Group_name      string
	Dialed_cnt      string
	Company_name    string
	Call_time       string
	Peer_num_loc    string
	Location        string
	Peer_number_cnt string
}
type CellBehaviorItem struct {
	Phone_num string
	Behavior  []Behavior
}
type Behavior struct {
	Cell_operator_zh string
	Net_flow         string
	Dialed_cnt       string
	Rechange_cnt     string
	Call_time        string
	Cell_operator    string
	Cell_phone_num   string
	Sms_cnt          string
	Cell_loc         string
	Dialed_time      string
	Call_cnt         string
	Dial_cnt         string
	Dial_time        string
	Total_amount     string
	Cell_mth         string
	Rechange_amount  string
}
type ConsumptionDetailItem struct {
	Item         ConsumptionDetailItemItem
	App_point    string
	App_point_zh string
}
type ConsumptionDetailItemItem struct {
	Item_6m     string
	Item_3m     string
	Avg_item_3m string
	Item_1m     string
	Avg_item_6m string
}
type BasicCheckItem struct {
	Result     string
	Comment    string
	Check_item string
}
type ContactRegionItem struct {
	Region_list []RegionListItem
	Key         string
	Desc        string
}
type RegionListItem struct {
	Region_dialed_time     string
	Region_dial_time       string
	Region_dial_cnt_pct    string
	Region_dial_time_pct   string
	Region_call_cnt        string
	Region_call_time       string
	Region_avg_dialed_time string
	Region_dialed_cnt      string
	Region_loc             string
	Region_dial_cnt        string
	Region_dialed_time_pct string
	Region_avg_dial_time   string
	Region_uniq_num_cnt    string
	Region_dialed_cnt_pct  string
}
type CallDurationDetailItem struct {
	Duration_list []DurationListItem
	Key           string
	Desc          string
}
type DurationListItem struct {
	Time_step    string
	Item         DurationListItemItem
	Time_step_zh string
}
type DurationListItemItem struct {
	Dialed_time        string
	Total_cnt          string
	Dial_cnt           string
	Dial_time          string
	Dialed_cnt         string
	Latest_call_time   string
	Farthest_call_time string
	Uniq_num_cnt       string
	Total_time         string
}
type CallTimeDetailItem struct {
	Item         CallTimeDetailItemItem
	App_point    string
	App_point_zh string
}
type CallTimeDetailItemItem struct {
	Item_6m     string
	Item_3m     string
	Avg_item_3m string
	Item_1m     string
	Avg_item_6m string
}
type MainServiceItem struct {
	Total_service_cnt string
	Service_details   []ServiceDetailsItem
	Group_name        string
	Company_name      string
	Service_num       string
}
type ServiceDetailsItem struct {
	Dialed_time   string
	Interact_mth  string
	Interact_time string
	Dial_cnt      string
	Dial_time     string
	Dialed_cnt    string
	Interact_cnt  string
}
type ActiveDegreeItem struct {
	Item         ActiveDegreeItemItem
	App_point    string
	App_point_zh string
}
type ActiveDegreeItemItem struct {
	Item_6m     string
	Item_3m     string
	Avg_item_3m string
	Item_1m     string
	Avg_item_6m string
}
type CallServiceAnalysis struct {
	Analysis_item  string
	Analysis_point AnalysisPoint
	Analysis_desc  string
}
type RoamDetailItem struct {
	Roam_day      string
	Roam_location string
}
type SmsContactDetailItem struct {
	Send_cnt_3m    string
	Receive_cnt_6m string
	Send_cnt_6m    string
	Sms_cnt_1m     string
	Receive_cnt_3m string
	Sms_cnt_3m     string
	Peer_num       string
	Sms_cnt_1w     string
	Sms_cnt_6m     string
}
type CallFamilyDetailItem struct {
	Item         CallFamilyDetailItemItem
	App_point    string
	App_point_zh string
}
type CallFamilyDetailItemItem struct {
	Item_6m     string
	Item_3m     string
	Avg_item_3m string
	Item_1m     string
	Avg_item_6m string
}
type CallContactDetailItem struct {
	Dialed_time_3m        string
	Call_if_whole_day_3m  string
	City                  string
	Call_cnt_morning_3m   string
	Call_cnt_night_3m     string
	Call_cnt_6m           string
	Call_cnt_evening_6m   string
	Call_cnt_weekday_3m   string
	Call_cnt_holiday_6m   string
	Peer_num              string
	Call_cnt_weekend_3m   string
	Call_cnt_afternoon_6m string
	P_relation            string
	Trans_start           string
	Call_time_3m          string
	Dial_cnt_3m           string
	Dial_time_3m          string
	Call_cnt_morning_6m   string
	Call_cnt_weekend_6m   string
	Call_cnt_holiday_3m   string
	Group_name            string
	Call_cnt_evening_3m   string
	Dialed_time_6m        string
	Call_if_whole_day_6m  string
	Call_cnt_weekday_6m   string
	Dialed_cnt_3m         string
	Call_cnt_night_6m     string
	Call_cnt_noon_3m      string
	Call_cnt_1m           string
	Dialed_cnt_6m         string
	Call_cnt_afternoon_3m string
	Company_name          string
	Call_cnt_3m           string
	Call_cnt_noon_6m      string
	Dial_time_6m          string
	Dial_cnt_6m           string
	Call_cnt_1w           string
	Call_time_6m          string
	Trans_end             string
}
type RoamAnalysisItem struct {
	Max_roam_day_cnt_6m  string
	Roam_day_cnt_6m      string
	Max_roam_day_cnt_3m  string
	Continue_roam_cnt_3m string
	Roam_location        string
	Continue_roam_cnt_6m string
	Roam_day_cnt_3m      string
}
type BehaviorCheckItem struct {
	Result         string
	Score          string
	Evidence       string
	Check_point_cn string
	Check_point    string
}
type ReData struct {
	User_basic            []KVItem
	Friend_circle         FriendCircle
	Call_risk_analysis    []CallRiskAnalysisItem
	Cell_behavior         []CellBehaviorItem
	Consumption_detail    []ConsumptionDetailItem
	Basic_check_items     []BasicCheckItem
	Contact_region        []ContactRegionItem
	Call_duration_detail  []CallDurationDetailItem
	Application_check     interface{}
	Call_time_detail      []CallTimeDetailItem
	Trip_info             interface{}
	Main_service          []MainServiceItem
	Active_degree         []ActiveDegreeItem
	Call_service_analysis []CallServiceAnalysis
	User_info_check       interface{}
	Cell_phone            []KVItem
	Roam_detail           []RoamDetailItem
	Sms_contact_detail    []SmsContactDetailItem
	Call_family_detail    []CallFamilyDetailItem
	Call_contact_detail   []CallContactDetailItem
	Roam_analysis         []RoamAnalysisItem
	Collection_contact    interface{}
	Report                []KVItem
	Behavior_check        []BehaviorCheckItem
}
type ApplicationCheckItem struct {
	App_point    string
	Check_points ApplicationCheckCheckPoints
}
type ApplicationCheckCheckPoints struct {
	Relationship  string
	Key_value     string
	Contact_name  string
	Check_mobile  string
	Check_xiaohao string
}
type TripInfoItem struct {
	Trip_dest       string
	Trip_start_time string
	Trip_end_time   string
	Trip_leave      string
	Trip_type       string
}
type UserInfoCheck struct {
	Check_search_info CheckSearchInfo
	Check_black_info  []CheckBlackInfo
}
type CheckSearchInfo struct {
	Searched_org_cnt         string
	Searched_org_type        string
	Idcard_with_other_names  string
	Idcard_with_other_phones string
	Phone_with_other_names   string
	Phone_with_other_idcards string
	Register_org_cnt         string
	Register_org_type        string
	Arised_open_web          string
}
type CheckBlackInfo struct {
	Phone_gray_score              string
	Contacts_class1_blacklist_cnt string
	Contacts_class2_blacklist_cnt string
	Contacts_class1_cnt           string
	Contacts_router_cnt           string
	Contacts_router_ratio         string
}
