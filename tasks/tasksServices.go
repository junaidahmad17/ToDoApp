package tasks

import (
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"
	"todoapp/email"
	"todoapp/users"

	"github.com/gin-gonic/gin"
	"github.com/robfig/cron"
)

var attachfolder string = os.Getenv("DBADD") + "attachments"

/////////////////////////////////////////////////////////////////////////////////////////////////
func GetUid(c *gin.Context) int {
	y, _ := c.Get("client")
	x, _ := strconv.Atoi(y.(string))
	return x
}


// Listing All Tasks
func GetTasks(c *gin.Context) {

	var task []Task
	if e := DB.Where(&Task{Uid: GetUid(c)}).Find(&task).Error; e != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": e.Error()})
		return
	}
	c.JSON(http.StatusOK, task)

}

// Creating a New Task
func CreateTask(c *gin.Context) {
	DB.AutoMigrate(&Task{})

	var input Task
	r, _ := c.Cookie("Token")
	if r == "" {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}
	c.BindJSON(&input)

	input.Uid = GetUid(c)

	var task Task
	if e := DB.Where("Title=? AND Uid=?", input.Title, GetUid(c)).First(&task).Error; e == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "task already exists!"})
		return
	}
	DB.Create(&input)
	
	c.JSON(http.StatusOK, gin.H{"msg": "task added!"})

}

// Task Deletion
func DeleteTask(c *gin.Context) {

	var task Task
	if e := DB.Where("ID=? AND Uid=?", c.Param("id"), GetUid(c)).First(&task).Error; e != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "task does not exist!"})
		return
	}

	DB.Delete(&task)
	c.JSON(http.StatusOK, gin.H{"msg": "task deleted successfully!"})
}

// Editing a Tasking Using ID as Key
func EditTask(c *gin.Context) {

	var task Task
	if e := DB.Where("ID=? AND Uid=?", c.Param("id"), GetUid(c)).First(&task).Error; e != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "task does not exist!"})
		return
	}

	var input UpdateTask
	c.BindJSON(&input)

	var check Task
	if e := DB.Where("Title=? AND Uid=?", input.Title, GetUid(c)).First(&check).Error; e == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "task already exists!"})
		return
	}

	DB.Model(&task).Updates(Task{Title: input.Title, Description: input.Description, Com_status: input.Com_status})
	c.JSON(http.StatusOK, gin.H{"msg": "task modified successfully!"})
}

// Clearing Whole DB
func DeleteAll(c *gin.Context) {
	var task0 []Task
	DB.Find(&task0)
	DB.Delete(&task0)
	c.JSON(http.StatusOK, gin.H{"msg": "all entries deleted!"})
}

////////////////////////////           Reports                //////////////////////////////////////////////
// Listing Basic Stats
func CountTask(c *gin.Context) {
	var task []Task
	// Total Tasks
	out := DB.Where("Uid=?", GetUid(c)).Find(&task)
	c.JSON(http.StatusOK, gin.H{"total tasks": out.RowsAffected})
	// Completed Tasks
	out = DB.Where(&Task{Com_status: true, Uid: GetUid(c)}).Find(&task)
	c.JSON(http.StatusOK, gin.H{"completed tasks": out.RowsAffected})
	// Remaining Tasks
	out = DB.Where(&Task{}, "Com_status").Find(&task)
	c.JSON(http.StatusOK, gin.H{"remaining tasks": out.RowsAffected})
}

// Incomplete Tasks Past Due Dates
func MissedTasks(c *gin.Context) {

	t := time.Now()
	count := 0

	rows, _ := DB.Model(&Task{}).Where("Com_status = ? AND Uid =?", false, GetUid(c)).Rows()
	var task Task
	for rows.Next() {
		DB.ScanRows(rows, &task)

		if task.Due_DT.Sub(t) < 0 {
			count = count + 1
		}
	}

	c.JSON(http.StatusOK, gin.H{"missed deadlines": count})
}

func AttachFile(c *gin.Context) {
	var task Task
	if e := DB.Where("id=? AND Uid=?", c.Param("id"), GetUid(c)).First(&task).Error; e != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "task does not exist!"})
		return
	}
	file, fh, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error})
	}

	filename := "u" + strconv.Itoa(GetUid(c)) + "tsk" + c.Param("id") + fh.Filename
	path := attachfolder + "\\" + filename

	_, err = os.Stat(path)
	if os.IsExist(err) {
		os.Remove(path)
	}
	nfile, err := os.Create(path)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error})
		return
	}
	defer nfile.Close()
	content, err := ioutil.ReadAll(file)
	if err != nil {

		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error})
		return
	}
	defer file.Close()
	nfile.Write(content)
	task.Attachment = fh.Filename
	DB.Save(&task)
	c.JSON(http.StatusOK, "File uploaded successfully")

}

func DeleteFile(c *gin.Context) {
	var task Task
	if e := DB.Where("id=? AND Uid=?", c.Param("id"), GetUid(c)).First(&task).Error; e != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "task does not exist!"})
		return
	}
	if task.Attachment == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": "no attachment found!"})
		return
	}
	filename := "u" + strconv.Itoa(GetUid(c)) + "tsk" + c.Param("id") + task.Attachment
	path := attachfolder + "\\" + filename
	os.Remove(path)
	task.Attachment = ""
	DB.Save(&task)
	c.JSON(http.StatusOK, gin.H{"error": "attachment deleted successfully!"})
}

func DownloadFile(c *gin.Context) {
	var task Task
	if e := DB.Where("id=? AND Uid=?", c.Param("id"), GetUid(c)).First(&task).Error; e != nil {
		c.JSON(http.StatusOK, gin.H{"error": "task does not exist!"})
		return
	}
	if task.Attachment == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": "the task does not have any attachment"})
		return
	}
	filename := "u" + strconv.Itoa(GetUid(c)) + "tsk" + c.Param("id") + task.Attachment
	path := attachfolder + "\\" + filename
	file, err := os.Open(path)
	if err != nil {
		c.JSON(http.StatusNotImplemented, gin.H{"error": "unable to open file!"})
		return
	}
	defer file.Close()
	mkfile := make([]byte, 512)
	file.Read(mkfile)
	fileType := http.DetectContentType(mkfile)
	fileInfo, _ := file.Stat()
	fileSize := fileInfo.Size()

	c.Header("Attachment name: ", task.Attachment)
	c.Header("type: ", fileType)
	fs := strconv.FormatInt(fileSize, 10)
	c.Header("size: ", fs)

	file.Seek(0, 0)
	io.Copy(c.Writer, file)

}

func SimilarTasks(c *gin.Context) {

	rows, _ := DB.Model(&Task{}).Where("Uid = ?", GetUid(c)).Order("Description").Rows()
	defer rows.Close()
	var description string
	count := 2
	flag := false
	c.JSON(http.StatusOK, "Set 1\n")
	for rows.Next() {
		var task Task
		DB.ScanRows(rows, &task)

		if flag && description != task.Description {
			c.JSON(http.StatusOK, "Set "+strconv.FormatInt(int64(count), 10))
			count = count + 1
		}
		c.JSON(http.StatusOK, task.Description)
		description = task.Description
		flag = true
	}
}

func Remind() {
	cr := cron.New() 
	cr.AddFunc("00 00 12 * * *", func() {
		rows, _ := DB.Model(&Task{}).Order("Due_DT").Rows()
		defer rows.Close()
		t := time.Now()
		for rows.Next() {
			var task Task
			var user users.User

			DB.ScanRows(rows, &task)
			
			if task.Due_DT.Day()==t.Day() && task.Due_DT.Month()==t.Month() && task.Due_DT.Year()==t.Year() {
				users.UDB.Where(&users.User{ID:uint(task.Uid)}).First(&user)
				str :=  "Dear "+user.Username +",\n\nYou have a deadline coming up today!\n\nTitle: "+task.Title
				str = str+"\n\nRegards,\nToDo Team"
				email.SendEmail(user.Email, "Deadline Reminder", str) 
			}
		}

	}) 
	cr.Start()
}
