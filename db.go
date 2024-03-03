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
func incrementRevertCount(vmName string) int {
    var revertCount RevertCount 
    teamName := vmName[:4]

    res := db.Where("vm_name = ?", vmName).First(&revertCount)
    if res.Error != nil {
        revertCount = RevertCount{VMName: vmName, TeamName: teamName, Count: 1}
        db.Create(&revertCount)
        return 1
    }

    revertCount.Count++
    db.Save(&revertCount)

    return revertCount.Count
}

// Gets the revert counts for a Team
func getTeamRevertCount(teamName string) ([]RevertCount, error) {
    var revertCounts []RevertCount
    tx := db.Where("team_name = ?", teamName).Find(&revertCounts)
    if tx.Error != nil {
        return nil, tx.Error
    }
    return revertCounts, nil
}

// Gets all revert counts
func getAllRevertCounts() []RevertCount {
    var revertCounts []RevertCount
    db.Find(&revertCounts)
    return revertCounts
}
