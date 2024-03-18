package server

import (
	"encoding/json"
	"hippias-fiber/internal/models"
	"log"

	"github.com/gofiber/fiber/v2"
)

func (s *Server) GetDiscussionMgmtDetails(c *fiber.Ctx) error {
	discussionID := c.Params("id")

	var discussionResult json.RawMessage
	err := s.sb.DB.From("discussions").
		Select("*").
		Single().
		Eq("id", discussionID).
		Execute(&discussionResult)
	if err != nil {
		log.Printf("Error querying discussion: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Message: err.Error()})
	}

	var discussion models.Discussion
	err = json.Unmarshal(discussionResult, &discussion)
	if err != nil {
		log.Printf("Error unmarshaling discussion: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Message: err.Error()})
	}

	var participantsResult json.RawMessage
	err = s.sb.DB.From("course_participants").
		Select("*").
		Eq("course_id", string(discussion.CourseID)).
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

	var readingsResult json.RawMessage
	err = s.sb.DB.From("readings").
		Select("*").
		Eq("discussion_id", discussionID).
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

	var readingDtos []models.ReadingDto
	for _, reading := range readings {

		var ratingsResult json.RawMessage
		err := s.sb.DB.From("reading_ratings").
			Select("*").
			Eq("reading_id", string(reading.ID)).
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

		readingDto := models.ReadingDto{
			Reading: reading,
			Ratings: ratings,
		}
		readingDtos = append(readingDtos, readingDto)
	}

	var attendanceResult json.RawMessage
	err = s.sb.DB.From("discussion_attendance").
		Select("*").
		Eq("discussion_id", discussionID).
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

	discussionMgmtDto := models.DiscussionMgmtDto{
		Discussion:   discussion,
		Participants: participantDtos,
		Readings:     readingDtos,
		Attendance:   attendance,
	}

	return c.JSON(discussionMgmtDto)
}
