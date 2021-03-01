package enum

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type UserType Enum
type UserTypeEnum struct {
	Free UserType
	Vip  UserType
}
type NestUserTypeEnum struct {
	First  UserTypeEnum
	Second UserTypeEnum
}

func TestAssign(t *testing.T) {
	userTypeEnum := &UserTypeEnum{}
	Assign(userTypeEnum)
	require.Equal(t, UserType(0), userTypeEnum.Free)
	require.Equal(t, UserType(1), userTypeEnum.Vip)
	nestUserTypeEnum := &NestUserTypeEnum{}
	Assign(nestUserTypeEnum)
	require.Equal(t, UserType(2), nestUserTypeEnum.First.Free)
	require.Equal(t, UserType(3), nestUserTypeEnum.First.Vip)
	require.Equal(t, UserType(4), nestUserTypeEnum.Second.Free)
	require.Equal(t, UserType(5), nestUserTypeEnum.Second.Vip)
}
