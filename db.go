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
func incrementRevertCount(teamID string, vmName string) {
    var revertCount RevertCount
    err := db.First(&revertCount)
    if err != nil {
        revertCount = RevertCount{TeamID: teamID, VMName: vmName, Count: 0}
        db.Create(&revertCount)
    }
    revertCount.Count++
    db.Save(&revertCount)
}

// Gets the revert counts for a Team
func getRevertCount(teamID string) ([]RevertCount, error) {
    var revertCounts []RevertCount
    err := db.Find(&revertCounts).Where("team_id = ?", teamID)
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
