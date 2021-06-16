package tasks

import (
	//"fmt"
	//"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)
var attachfolder string = "C:\\Users\\Junaid Ahmad (WORK)\\Desktop\\GO\\newtodo\\todoapp\\attachments"

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
		c.JSON(http.StatusBadRequest, gin.H{"Error": e.Error()})
		return 
	}
	c.JSON(http.StatusOK,task)
	
}

// Creating a New Task
func CreateTask(c *gin.Context) {
	DB.AutoMigrate(&Task{})
	
	var input Task
	r,_ := c.Cookie("Token")
	if r == "" {
		c.JSON(http.StatusForbidden,"Forbidden")
		return
	}
	c.BindJSON(&input)

	input.Uid = GetUid(c)

	var task Task 
	if e := DB.Where("Title=? AND Uid=?",input.Title,GetUid(c)).First(&task).Error; e == nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Task already exists!"})
		return 
	}
	DB.Create(&input)

	c.JSON(http.StatusOK, "Task Added!")
	
}

// Task Deletion
func DeleteTask(c *gin.Context) {
	
	var task Task
	if e := DB.Where("ID=? AND Uid=?",c.Param("id"),GetUid(c)).First(&task).Error; e != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Task does not exist!"})
		return 
	}

	DB.Delete(&task)
	c.JSON(http.StatusOK, "Task Deleted Successfully!")
}

// Editing a Tasking Using ID as Key 
func EditTask(c *gin.Context) {

	var task Task
	if e := DB.Where("ID=? AND Uid=?",c.Param("id"),GetUid(c)).First(&task).Error; e != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Task does not exist!"})
		return 
	}
	
	var input UpdateTask
	c.BindJSON(&input)

	var check Task
	if e := DB.Where("Title=? AND Uid=?",input.Title,GetUid(c)).First(&check).Error; e == nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Task already exists!"})
		return 
	}

	DB.Model(&task).Updates(Task{Title: input.Title, Description: input.Description, Com_status: input.Com_status})
	c.JSON(http.StatusOK, "Task Modified Successfully!")
}

// Clearing Whole DB
func DeleteAll(c *gin.Context) {
	var task0 []Task
	DB.Find(&task0)
	DB.Delete(&task0)
	c.JSON(http.StatusOK, "All Entries Deleted!")
}

////////////////////////////           Reports                //////////////////////////////////////////////
// Listing Basic Stats
func CountTask(c *gin.Context) {
	var task []Task
	// Total Tasks 
	out := DB.Where("Uid=?",GetUid(c)).Find(&task)
	c.JSON(http.StatusOK, gin.H{"Total Tasks": out.RowsAffected})
	// Completed Tasks 
	out = DB.Where(&Task{Com_status: true,Uid: GetUid(c) }).Find(&task)
	c.JSON(http.StatusOK, gin.H{"Completed Tasks": out.RowsAffected})
	// Remaining Tasks
	out = DB.Where(&Task{},"Com_status").Find(&task)
	c.JSON(http.StatusOK, gin.H{"Remaining Tasks": out.RowsAffected})
}

// Incomplete Tasks Past Due Dates
func MissedTasks(c *gin.Context) {

	t := time.Now()
	count := 0

	rows, _ := DB.Model(&Task{}).Where("Com_status = ? AND Uid =?", false,GetUid(c)).Rows()
	var task Task
	for rows.Next() {
		DB.ScanRows(rows,&task)
		
		if task.Due_DT.Sub(t) < 0 {
			count = count + 1
		}
	}

	c.JSON(http.StatusOK, gin.H{"Missed Deadlines":count})
}


func AttachFile(c *gin.Context) {
	var task Task
	if e := DB.Where("id=? AND Uid=?",c.Param("id"),GetUid(c)).First(&task).Error; e != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Task does not exist!"})
		return 
	}
	file, fh, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusNotFound,err.Error)
	}
	
	filename := "u"+strconv.Itoa(GetUid(c))+"tsk"+c.Param("id")+fh.Filename
	path := attachfolder+"\\"+filename
	
	_, err = os.Stat(path)
	if os.IsExist(err) {
		os.Remove(path)
	}
	nfile, err := os.Create(path)
	if err != nil {
		c.JSON(http.StatusBadRequest,err.Error())
		return
	}
	defer nfile.Close()
	content, err := ioutil.ReadAll(file)
	if err != nil {
			
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	defer file.Close()
	nfile.Write(content)
	task.Attachment = fh.Filename
	DB.Save(&task)
	c.JSON(http.StatusOK,"File uploaded successfully")

}

func DeleteFile(c *gin.Context) {
	var task Task
	if e := DB.Where("id=? AND Uid=?",c.Param("id"),GetUid(c)).First(&task).Error; e != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Task does not exist!"})
		return 
	}
	if task.Attachment == "" {
		c.JSON(http.StatusNotFound, "No Attachment Found!")
		return 
	}
	filename := "u"+strconv.Itoa(GetUid(c))+"tsk"+c.Param("id")+task.Attachment
	path := attachfolder+"\\"+filename
	os.Remove(path)
	task.Attachment = ""
	DB.Save(&task)
	c.JSON(http.StatusOK, "Attachment Deleted Successfully!")
}

func DownloadFile(c *gin.Context) {
	var task Task
	if e := DB.Where("id=? AND Uid=?",c.Param("id"),GetUid(c)).First(&task).Error; e != nil {
		c.JSON(http.StatusOK, "Task does not exist!")
		return 
	}
	if task.Attachment == "" {
		c.JSON(http.StatusNotFound, "The task does not have any attachment")
		return 
	}
	filename := "u"+strconv.Itoa(GetUid(c))+"tsk"+c.Param("id")+task.Attachment
	path := attachfolder+"\\"+filename
	file, err := os.Open(path)
	if err!= nil {
		c.JSON(http.StatusNotImplemented, "Unable to open file!")
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

	file.Seek(0,0)
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
			c.JSON(http.StatusOK, "Set "+strconv.FormatInt(int64(count),10))
			count = count+1
		}
		c.JSON(http.StatusOK, task.Description)
		description = task.Description
		flag = true
	}
}