/*  Validation.go
*
* @Author:             Audy Vaksi <vaksipranata@gmail.com>
* @Date:               September 07, 2018
* @Last Modified by:   @vaksi
* @Last Modified time: 07/09/18 00:52
 */

package helpers

import (
	"github.com/pkg/errors"
	"github.com/vaksi/foreign_currency/constants"
)

func ValidateInValueCurrency(value interface{}) error {
	s, _ := value.(string)
	if constants.Currency(s).Code() == "" {
		return errors.Errorf("Value %s is not valid in Currency", s)
	}
	return nil
}
