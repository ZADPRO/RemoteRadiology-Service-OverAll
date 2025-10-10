package service

import (
	hashdb "AuthenticationService/internal/Helper/HashDB"
	logger "AuthenticationService/internal/Helper/Logger"
	helper "AuthenticationService/internal/Helper/ViewFile"
	model "AuthenticationService/internal/Model/ProfileService"
	query "AuthenticationService/query/ProfileService"
	"log"
	"strings"

	"gorm.io/gorm"

)

func isS3URL(path string) bool {
	log.Printf("IS S3 URL checkin : %v", path)
	return strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://")
}

func GetAllCoDoctorDataService(db *gorm.DB, reqVal model.GetReceptionistReq) []model.GetAllRadiologist {
	log := logger.InitLogger()

	var RadiologistData []model.GetAllRadiologist

	err := db.Raw(query.GetListofCoDoctorSQL, 8, reqVal.ScanID).Scan(&RadiologistData).Error
	if err != nil {
		log.Printf("ERROR: Failed to fetch scan centers: %v", err)
		return []model.GetAllRadiologist{}
	}

	for i, tech := range RadiologistData {
		RadiologistData[i].FirstName = hashdb.Decrypt(tech.FirstName)
		RadiologistData[i].LastName = hashdb.Decrypt(tech.LastName)
	}

	return RadiologistData
}

func GetDoctorCoDataService(db *gorm.DB, reqVal model.GetOneReceptionistReq, idValue int) []model.GetCoDoctorOne {
	log := logger.InitLogger()
	var RadiologistData []model.GetCoDoctorOne

	UserId := reqVal.UserId
	ScanCenterId := reqVal.ScanID

	if reqVal.UserId == 0 {
		UserId = idValue
		var MappingData []model.Mapping
		if err := db.Raw(query.IdentifyScanCenterMapping, idValue).Scan(&MappingData).Error; err != nil {
			log.Printf("ERROR: Failed to fetch scan centers: %v", err)
			return []model.GetCoDoctorOne{}
		}
		if len(MappingData) > 0 {
			ScanCenterId = MappingData[0].SCId
		} else {
			ScanCenterId = 0
		}
	}

	if err := db.Raw(query.GetListofCoDoctorOneSQL, UserId, ScanCenterId).Scan(&RadiologistData).Error; err != nil {
		log.Printf("ERROR: Failed to fetch co-doctors: %v", err)
		return []model.GetCoDoctorOne{}
	}

	for i, tech := range RadiologistData {
		// 🔹 Decrypt all fields
		RadiologistData[i].FirstName = hashdb.Decrypt(tech.FirstName)
		RadiologistData[i].LastName = hashdb.Decrypt(tech.LastName)
		RadiologistData[i].ProfileImg = hashdb.Decrypt(tech.ProfileImg)
		RadiologistData[i].DOB = hashdb.Decrypt(tech.DOB)
		RadiologistData[i].SocialSecurityNo = hashdb.Decrypt(tech.SocialSecurityNo)
		RadiologistData[i].DriversLicenseNo = hashdb.Decrypt(tech.DriversLicenseNo)
		RadiologistData[i].Specialization = hashdb.Decrypt(tech.Specialization)
		RadiologistData[i].DigitalSignature = hashdb.Decrypt(tech.DigitalSignature)
		RadiologistData[i].NPI = hashdb.Decrypt(tech.NPI)

		// ---------- PROFILE IMAGE ----------
		profilePath := RadiologistData[i].ProfileImg
		if len(profilePath) > 0 {
			if isS3URL(profilePath) {
				log.Printf("\n\nProfile path -> if condition %v", profilePath)
				// S3 URL → directly use it
				RadiologistData[i].ProfileImgFile = &model.FileData{Base64Data: profilePath, ContentType: "url"}
			} else {
				log.Printf("\n\nProfile path -> else condition %v", profilePath)
				// Local file fallback
				profileImgHelperData, err := helper.ViewFile("./Assets/Profile/" + profilePath)
				if err != nil {
					log.Errorf("\n\n\nFailed to read profile image: %v", err)
				} else if profileImgHelperData != nil {
					RadiologistData[i].ProfileImgFile = &model.FileData{
						Base64Data:  profileImgHelperData.Base64Data,
						ContentType: profileImgHelperData.ContentType,
					}
				}
			}
		}

		// ---------- DRIVER LICENSE ----------
		driverPath := RadiologistData[i].DriversLicenseNo
		if len(driverPath) > 0 {
			if isS3URL(driverPath) {
				RadiologistData[i].DriversLicenseFile = &model.FileData{Base64Data: driverPath, ContentType: "url"}
			} else {
				driverData, err := helper.ViewFile("./Assets/Files/" + driverPath)
				if err != nil {
					log.Errorf("Failed to read DrivingLicense file: %v", err)
				} else if driverData != nil {
					RadiologistData[i].DriversLicenseFile = &model.FileData{
						Base64Data:  driverData.Base64Data,
						ContentType: driverData.ContentType,
					}
				}
			}
		}

		// ---------- DIGITAL SIGNATURE ----------
		signPath := RadiologistData[i].DigitalSignature
		if len(signPath) > 0 {
			if isS3URL(signPath) {
				RadiologistData[i].DigitalSignatureFile = &model.FileData{Base64Data: signPath, ContentType: "url"}
			} else {
				signData, err := helper.ViewFile("./Assets/Profile/" + signPath)
				if err != nil {
					log.Errorf("Failed to read DigitalSignature: %v", err)
				} else if signData != nil {
					RadiologistData[i].DigitalSignatureFile = &model.FileData{
						Base64Data:  signData.Base64Data,
						ContentType: signData.ContentType,
					}
				}
			}
		}

		// ---------- LICENSE FILES ----------
		var licenseFiles []model.LicenseFilesModel
		if err := db.Raw(query.GetLicenseFilesSQL, reqVal.UserId).Scan(&licenseFiles).Error; err == nil {
			for _, lf := range licenseFiles {
				lf.LFileName = hashdb.Decrypt(lf.LFileName)
				if isS3URL(lf.LFileName) {
					lf.LFileData = &model.FileData{Base64Data: lf.LFileName, ContentType: "url"}
				} else {
					if data, err := helper.ViewFile("./Assets/Files/" + lf.LFileName); err == nil && data != nil {
						lf.LFileData = &model.FileData{
							Base64Data:  data.Base64Data,
							ContentType: data.ContentType,
						}
					}
				}
				RadiologistData[i].LicenseFiles = append(RadiologistData[i].LicenseFiles, lf)
			}
		}

		// ---------- MALPRACTICE FILES ----------
		var malpracticeFiles []model.MalpracticeModel
		if err := db.Raw(query.GetMalpracticeFilesSQL, reqVal.UserId).Scan(&malpracticeFiles).Error; err == nil {
			for _, mp := range malpracticeFiles {
				mp.MPFileName = hashdb.Decrypt(mp.MPFileName)
				if isS3URL(mp.MPFileName) {
					mp.MPFileData = &model.FileData{Base64Data: mp.MPFileName, ContentType: "url"}
				} else {
					if data, err := helper.ViewFile("./Assets/Files/" + mp.MPFileName); err == nil && data != nil {
						mp.MPFileData = &model.FileData{
							Base64Data:  data.Base64Data,
							ContentType: data.ContentType,
						}
					}
				}
				RadiologistData[i].MalpracticeInsuranceDetails = append(RadiologistData[i].MalpracticeInsuranceDetails, mp)
			}
		}
	}

	return RadiologistData
}
