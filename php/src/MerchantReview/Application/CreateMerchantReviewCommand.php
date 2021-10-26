<?php

declare(strict_types=1);

namespace App\MerchantReview\Application;

final class CreateMerchantReviewCommand
{
    public function __construct(
        private string $comment,
		private int $rating,
    ) {}

    public function getComment(): string
    {
        return $this->comment;
    }

    public function getRating(): int
    {
        return $this->rating;
    }
}
