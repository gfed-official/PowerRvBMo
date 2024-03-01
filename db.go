package main

import (
    "gorm.io/gorm"
    "gorm.io/driver/sqlite"
)

// Connects to Database
func dbConn() (db *gorm.DB) {
    db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})

    if err != nil {
        panic("failed to connect database")
    }

    return db
}

// Increments the revert count for a VM
func incrementRevertCount(vmName string) {
    var revertCount RevertCount 
    res := db.Model(&RevertCount{VMName: vmName}).First(&revertCount)
    if res.Error != nil {
        revertCount = RevertCount{VMName: vmName, Count: 1}
        db.Create(&revertCount)
        return
    }

    revertCount.Count++
    db.Save(&revertCount)
}

// Gets the revert counts for a Team
func getRevertCount(teamName string) ([]RevertCount, error) {
    var revertCounts []RevertCount
    err := db.Find(&revertCounts).Where("teamname = ?", teamName)
    if err != nil {
        return nil, err.Error
    }
    return revertCounts, nil
}

// Gets all revert counts
func getAllRevertCounts() []RevertCount {
    var revertCounts []RevertCount
    db.Find(&revertCounts)
    return revertCounts
}
