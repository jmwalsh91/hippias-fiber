package server

import (
	"encoding/json"
	"hippias-fiber/internal/models"
	"log"

	"github.com/gofiber/fiber/v2"
)

func (s *Server) getCourseManagementDetails(c *fiber.Ctx) error {
	courseID := c.Params("id")
	log.Printf("Fetching course management details for course %s", courseID)
	var courseResult json.RawMessage
	err := s.sb.DB.From("courses").
		Select("*").
		Single().
		Eq("id", courseID).
		Execute(&courseResult)
	if err != nil {
		log.Printf("Error querying course: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Message: err.Error()})
	}

	var course models.Course
	err = json.Unmarshal(courseResult, &course)
	if err != nil {
		log.Printf("Error unmarshaling course: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Message: err.Error()})
	}

	var discussionsResult json.RawMessage
	err = s.sb.DB.From("discussions").
		Select("*").
		Eq("course_id", courseID).
		Execute(&discussionsResult)
	if err != nil {
		log.Printf("Error querying discussions: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Message: err.Error()})
	}

	var discussions []models.Discussion
	err = json.Unmarshal(discussionsResult, &discussions)
	if err != nil {
		log.Printf("Error unmarshaling discussions: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Message: err.Error()})
	}

	var discussionDtos []models.DiscussionDto
	for _, discussion := range discussions {
		// Fetch readings for the discussion
		var readingsResult json.RawMessage
		err := s.sb.DB.From("readings").
			Select("*").
			Eq("discussion_id", string(discussion.ID)).
			Execute(&readingsResult)
		if err != nil {
			log.Printf("Error querying readings: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Message: err.Error()})
		}

		var readings []models.Reading
		err = json.Unmarshal(readingsResult, &readings)
		if err != nil {
			log.Printf("Error unmarshaling readings: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Message: err.Error()})
		}

		// Fetch reading ratings for the discussion
		var ratingsResult json.RawMessage
		err = s.sb.DB.From("reading_ratings").
			Select("*").
			Eq("discussion_id", string(discussion.ID)).
			Execute(&ratingsResult)
		if err != nil {
			log.Printf("Error querying reading ratings: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Message: err.Error()})
		}

		var ratings []models.ReadingRating
		err = json.Unmarshal(ratingsResult, &ratings)
		if err != nil {
			log.Printf("Error unmarshaling reading ratings: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Message: err.Error()})
		}

		// Fetch attendance for the discussion
		var attendanceResult json.RawMessage
		err = s.sb.DB.From("discussion_attendance").
			Select("*").
			Eq("discussion_id", string(discussion.ID)).
			Execute(&attendanceResult)
		if err != nil {
			log.Printf("Error querying discussion attendance: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Message: err.Error()})
		}

		var attendance []models.DiscussionAttendance
		err = json.Unmarshal(attendanceResult, &attendance)
		if err != nil {
			log.Printf("Error unmarshaling discussion attendance: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Message: err.Error()})
		}

		discussionDto := models.DiscussionDto{
			Discussion: discussion,
			Readings:   readings,
			Ratings:    ratings,
			Attendance: attendance,
		}
		discussionDtos = append(discussionDtos, discussionDto)
	}

	// Fetch course participants
	var participantsResult json.RawMessage
	err = s.sb.DB.From("course_participants").
		Select("*").
		Eq("course_id", courseID).
		Execute(&participantsResult)
	if err != nil {
		log.Printf("Error querying course participants: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Message: err.Error()})
	}

	var participants []models.CourseParticipant
	err = json.Unmarshal(participantsResult, &participants)
	if err != nil {
		log.Printf("Error unmarshaling course participants: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Message: err.Error()})
	}

	var participantDtos []models.CourseParticipantDto
	for _, participant := range participants {
		var userResult json.RawMessage
		err := s.sb.DB.From("users").
			Select("*").
			Single().
			Eq("id", string(participant.UserID)).
			Execute(&userResult)
		if err != nil {
			log.Printf("Error querying user: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Message: err.Error()})
		}

		var user models.User
		err = json.Unmarshal(userResult, &user)
		if err != nil {
			log.Printf("Error unmarshaling user: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Message: err.Error()})
		}

		participantDto := models.CourseParticipantDto{
			CourseParticipant: participant,
			User:              user,
		}
		participantDtos = append(participantDtos, participantDto)
	}

	courseMgmtDto := models.CourseMgmtDto{
		Course:       course,
		Discussions:  discussionDtos,
		Participants: participantDtos,
	}

	return c.JSON(courseMgmtDto)
}
