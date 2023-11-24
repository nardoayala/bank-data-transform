package main

import (
    "fmt"
    "os"
    "strconv"
    "strings"
    "github.com/xuri/excelize/v2"
    "github.com/atotto/clipboard"
)

type Config struct {
    SheetName string
    Columns []string
    StartRow int
}

func formatDate(date string) (string) {
    dateParts := strings.Split(date, "/")
    formattedDate := strings.Join([]string{
        dateParts[1], dateParts[0], dateParts[2],
    }, "/")
    return formattedDate
}

func formatNumber(num string) (string) {
    return strings.Replace(num, ",", "", -1)
}

func reverseSlice(s [][]string) {
    for i,j := 0, len(s)-1; i < j; i,j = i+1, j-1 {
        s[i], s[j] = s[j], s[i]
    }
}

func processRows(f *excelize.File, cfg Config, endRow int) ([][]string, error) {
    var allRows [][]string

    for i := cfg.StartRow; i <= endRow; i++ {
        var row []string
        for _, col := range cfg.Columns {
            cell, err := f.GetCellValue(
                cfg.SheetName, fmt.Sprintf("%s%d", col, i),
            )
            if err != nil {
                return nil, err
            }

            var value string = cell
            if col == "B" {
                value = formatDate(cell)
            }
            
            if col == "N" || col == "R" {
                value = formatNumber(cell)
            }

            row = append(row, value)
            if col == "B" {
                row = append(row, "")
            }
            if col == "N" {
                row = append(row, "Facebank")
            }
        }
        allRows = append(allRows, row)
    }
    return allRows, nil
}

func formatDataForClipboard(data [][]string) string {
    var builder strings.Builder
    for index, row := range data {
        builder.WriteString(strings.Join(row, "\t"))
        if index != len(data) -1 {
            builder.WriteString("\n")
        }
    }
    return builder.String()
}

func main() {
    cfg := Config{
        SheetName: "Page 1",
        Columns: []string{"B", "G", "N", "R"},
        StartRow: 25,
    }

    if len(os.Args) < 3 {
        fmt.Println("Usage: program <file> <endrow>")
        return
    }

    file := os.Args[1]
    if _, err := os.Stat(file); os.IsNotExist(err) {
        fmt.Println("File does not exist:", file)
        return
    }

    endRow, err := strconv.Atoi(os.Args[2])
    if err != nil {
        fmt.Println("Error parsing command line argument:", err)
        return
    }
    
    if endRow < cfg.StartRow {
        fmt.Println("End row must be greater than start row")
        return
    }

    f, err := excelize.OpenFile(file)
    if err != nil {
        fmt.Println("Error opening file:", err)
        return
    }

    defer func() {
        if err := f.Close(); err != nil {
            fmt.Println(err)
        }
    }()

    allRows, err := processRows(f, cfg, endRow)
    if err != nil {
        fmt.Println("Error processing rows:", err)
        return
    }

    reverseSlice(allRows)

    clipboardData := formatDataForClipboard(allRows)
    if err := clipboard.WriteAll(clipboardData); err != nil {
        fmt.Println("Error copying to clipboard:", err)
        return
    }
}
