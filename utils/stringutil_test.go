package utils

import (
	"github.com/chen-keinan/go-command-eval/common"
	"github.com/stretchr/testify/assert"
	"testing"
)

//Test_CheckType_Permission_OK test
func Test_CheckType_Permission_OK(t *testing.T) {
	evalExpr := "${0} <= 644"
	bench := UpdateCmdExprParam
	ti := bench([]string{"700"}, evalExpr)
	assert.Equal(t, ti, "700 <= 644")
}

//Test_CheckType_Owner_OK test
func Test_CheckType_Owner_OK(t *testing.T) {
	evalExpr := "'${0}' == 'root:root';"
	bench := UpdateCmdExprParam
	ti := bench([]string{"root:root"}, evalExpr)
	assert.Equal(t, ti, "'root:root' == 'root:root'")
}

//Test_CheckType_ProcessParam_OK test
func Test_CheckType_ProcessParam_OK(t *testing.T) {
	evalExpr := "'${0}' == 'false';"
	bench := UpdateCmdExprParam
	ti := bench([]string{"false"}, evalExpr)
	assert.Equal(t, ti, "'false' == 'false'")
}

//Test_CheckType_Multi_ProcessParam_OK test
func Test_CheckType_Multi_ProcessParam_OK(t *testing.T) {
	evalExpr := "'RBAC' IN (${0});"
	bench := UpdateCmdExprParam
	ti := bench([]string{"RBAC,bbb"}, evalExpr)

	assert.Equal(t, ti, "'RBAC' IN ('RBAC','bbb')")
}

//Test_CheckType_Multi_ProcessParam_OK test
func Test_CheckType_Multi_ProcessParam_RexOK(t *testing.T) {
	evalExpr := "'RBAC' IN (${0});"
	bench := UpdateCmdExprParam
	ti := bench([]string{common.GrepRegex}, evalExpr)
	assert.Equal(t, ti, "'RBAC' == ''")
}

//Test_CheckType_Owner_OK test
func Test_CheckType_Regex_OK(t *testing.T) {
	evalExpr := "'${0}' == 'root:root';"
	bench := UpdateCmdExprParam
	ti := bench([]string{common.GrepRegex}, evalExpr)
	assert.Equal(t, ti, "'' == 'root:root'")
}

//Test_CheckType_Regex_MultiParamType test
func Test_CheckType_Regex_MultiParamType(t *testing.T) {
	evalExpr := "'${0}' != 'root:root'; && 'root:root' IN (${0});"
	bench := UpdateCmdExprParam
	ti := bench([]string{"root:root"}, evalExpr)
	assert.Equal(t, ti, "'root:root' != 'root:root' && 'root:root' == 'root:root'")
}

//Test_CheckType_Regex_MultiParamTypeManyValues test
func Test_CheckType_Regex_MultiParamTypeManyValues(t *testing.T) {
	evalExpr := "'${0}' != 'root:root'; && 'root:root' IN (${0});"
	bench := UpdateCmdExprParam
	ti := bench([]string{"root:root,abc"}, evalExpr)
	assert.Equal(t, ti, "'root:root,abc' != 'root:root' && 'root:root' IN ('root:root','abc')")
}

//Test_CheckType_Regex_DiffParamTypeManyValues test
func Test_CheckType_Regex_DiffParamTypeManyValues(t *testing.T) {
	evalExpr := "'${1}' == 'kkk'; && '${0}' != 'root:root'; && 'root:root' IN (${0});"
	bench := UpdateCmdExprParam
	ti := bench([]string{"root:root,abc", "kkk"}, evalExpr)
	assert.Equal(t, ti, "'kkk' == 'kkk' && 'root:root,abc' != 'root:root' && 'root:root' IN ('root:root','abc')")
}

func Test_ExcludeAuditTest(t *testing.T) {
	et := ExcludeAuditTest([]string{"1.2.4"}, "1.2.5")
	assert.True(t, et)
	et = ExcludeAuditTest([]string{"1.2.4"}, "1.2.4")
	assert.False(t, et)
	et = ExcludeAuditTest([]string{}, "1.2.4")
	assert.False(t, et)
}

//Test_GetSpecificTestsToExecute test
func Test_GetSpecificTestsToExecute(t *testing.T) {
	l := GetAuditTestsList("i", "i=1.2.3,1.4.5")
	assert.Equal(t, l[0], "1.2.3")
	assert.Equal(t, l[1], "1.4.5")
	l = GetAuditTestsList("e", "")
	assert.Equal(t, l[0], "")
}

//Test_RemoveNewLineSuffix test
func Test_RemoveNewLineSuffix(t *testing.T) {
	s := RemoveNewLineSuffix("abc\n")
	assert.Equal(t, s, "abc")
	s = RemoveNewLineSuffix("abc\n134")
	assert.Equal(t, s, "abc\n134")
	s = RemoveNewLineSuffix("abc")
	assert.Equal(t, s, "abc")
}

//Test_AddNewLineToNonEmptyStr test
func Test_AddNewLineToNonEmptyStr(t *testing.T) {
	k := AddNewLineToNonEmptyStr("abc")
	assert.Equal(t, k, "abc\n")
	k = AddNewLineToNonEmptyStr("\n")
	assert.Equal(t, k, "\n")
	k = AddNewLineToNonEmptyStr("abc\n")
	assert.Equal(t, k, "abc\n")
}

//Test_ValidParam test
func Test_ValidParam(t *testing.T) {
	match, num := ValidParamData("aaaaa ${2}bbbb ")
	assert.True(t, match)
	assert.Equal(t, "2", num)

}

func Test_ReadPolicyExpr(t *testing.T) {
	evalExpr := "'${0}' != '';&& [${0} MATCH no_permission.policy QUERY example.policy_eval RETURN allow]"
	policy, err := ReadPolicyExpr(evalExpr)
	assert.NoError(t, err)
	assert.Equal(t, policy.PolicyName, "no_permission.policy")
	assert.Equal(t, policy.PolicyQueryParam, "example.policy_eval")
	assert.Equal(t, policy.EvalParamNum, 0)
}

func Test_ReadPolicyExprwithReturn(t *testing.T) {
	evalExpr := "'${0}' != '';&& [${0} MATCH no_permission.policy QUERY example.deny RETURN allow_policy,namespace]"
	policy, err := ReadPolicyExpr(evalExpr)
	assert.NoError(t, err)
	assert.Equal(t, policy.PolicyName, "no_permission.policy")
	assert.Equal(t, policy.PolicyQueryParam, "example.deny")
	assert.Equal(t, policy.EvalParamNum, 0)
	assert.Equal(t, len(policy.ReturnKeys), 2)
}

func Test_GetPolicyExpr(t *testing.T) {
	evalExpr := "'${0}' == 'root:root'; && [${0} MATCH no_deny.policy]"
	policyExpr := GetPolicyExpr(evalExpr)
	assert.Equal(t, policyExpr, "[${0} MATCH no_deny.policy]")
}

func Test_MatchPolicySingleReturn(t *testing.T) {
	policyValues := true
	returnData := []string{"match"}
	mpr := MatchPolicy(policyValues, returnData)
	assert.Equal(t, mpr.ReturnValues["match"], "true")
}

func Test_MatchPolicySingleMulti(t *testing.T) {
	policyValues := map[string]interface{}{"match": true}
	returnData := []string{"match"}
	mpr := MatchPolicy(policyValues, returnData)
	assert.Equal(t, mpr.ReturnValues["match"], "true")
}
